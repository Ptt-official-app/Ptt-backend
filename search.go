package main

import (
	"context"
	"math"
	"strings"

	apipb "github.com/Ptt-official-app/Ptt-backend/internal/proto/api"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
	"github.com/Ptt-official-app/go-bbs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Search reads the complete board index because the current go-bbs connector
// does not apply offset and length. Filtering and pagination are consequently
// performed after filtering so TotalPosts is the total number of matches.
func (s *server) Search(ctx context.Context, req *apipb.SearchRequest) (*apipb.SearchReply, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "search request is required")
	}
	boardID, err := resolveSearchBoard(ctx, s.usecase, req.GetRef())
	if err != nil {
		return nil, err
	}
	if req.GetLength() < 0 {
		return nil, status.Error(codes.InvalidArgument, "length must not be negative")
	}

	cond := &usecase.ArticleSearchCond{}
	for _, filter := range req.GetFilter() {
		if filter == nil {
			return nil, status.Error(codes.InvalidArgument, "search filter must not be nil")
		}
		switch filter.GetType() {
		case apipb.SearchFilter_TYPE_TITLE:
			cond.Title = filter.GetStringData()
		case apipb.SearchFilter_TYPE_EXACT_TITLE:
			cond.ExactTitle = filter.GetStringData()
		case apipb.SearchFilter_TYPE_AUTHOR:
			cond.Author = filter.GetStringData()
		case apipb.SearchFilter_TYPE_RECOMMEND:
			// The client uses positive thresholds for good recommendations and
			// negative thresholds for bad recommendations.
			if filter.GetNumberData() >= 0 {
				cond.RecommendCountGreaterEqual = int(filter.GetNumberData())
				cond.RecommendCountGreaterEqualIsSet = true
			} else {
				cond.RecommendCountLessEqual = int(filter.GetNumberData())
				cond.RecommendCountLessEqualIsSet = true
			}
		case apipb.SearchFilter_TYPE_MONEY, apipb.SearchFilter_TYPE_MARK, apipb.SearchFilter_TYPE_SOLVED:
			return nil, status.Errorf(codes.Unimplemented, "search filter %s is not supported", filter.GetType().String())
		case apipb.SearchFilter_TYPE_UNKNOWN:
			return nil, status.Error(codes.InvalidArgument, "unknown search filter")
		default:
			return nil, status.Errorf(codes.InvalidArgument, "invalid search filter type %d", filter.GetType())
		}
	}

	articles := s.usecase.GetBoardArticles(ctx, boardID, 0, ^uint(0), cond)
	start, end := searchPageBounds(len(articles), int64(req.GetOffset()), int64(req.GetLength()))
	posts := make([]*apipb.Post, 0, end-start)
	for i := start; i < end; i++ {
		posts = append(posts, articleToSearchPost(articles[i], i))
	}

	return &apipb.SearchReply{
		Posts:      posts,
		TotalPosts: int32(minInt(len(articles), math.MaxInt32)),
	}, nil
}

// resolveSearchBoard uses the same 1-based board-id convention as the
// existing board cache, while also accepting a case-insensitive board name.
func resolveSearchBoard(ctx context.Context, uc usecase.Usecase, ref *apipb.BoardRef) (string, error) {
	if ref == nil {
		return "", status.Error(codes.InvalidArgument, "board ref is required")
	}
	boards := uc.GetBoards(ctx, "")
	if ref.GetBid() > 0 {
		index := int64(ref.GetBid()) - 1
		if index < 0 || index >= int64(len(boards)) {
			return "", status.Errorf(codes.InvalidArgument, "invalid board bid: %d", ref.GetBid())
		}
		return boards[index].BoardID(), nil
	}
	if name := strings.TrimSpace(ref.GetName()); name != "" {
		for _, board := range boards {
			if strings.EqualFold(board.BoardID(), name) {
				return board.BoardID(), nil
			}
		}
		return "", status.Errorf(codes.InvalidArgument, "invalid board name: %s", name)
	}
	return "", status.Error(codes.InvalidArgument, "board ref must contain a bid or name")
}

// searchPageBounds preserves repository order. A negative offset counts from
// the end, which is the convention used by pttweb's search pagination.
func searchPageBounds(total int, offset, length int64) (int, int) {
	if length <= 0 || total == 0 {
		return 0, 0
	}
	var start int64
	if offset < 0 {
		start = int64(total) + offset
	} else {
		start = offset
	}
	if start < 0 {
		start = 0
	}
	if start >= int64(total) {
		return total, total
	}
	end := start + length
	if end > int64(total) {
		end = int64(total)
	}
	return int(start), int(end)
}

func articleToSearchPost(article bbs.ArticleRecord, index int) *apipb.Post {
	recommend := article.Recommend()
	if recommend > math.MaxInt32 {
		recommend = math.MaxInt32
	} else if recommend < math.MinInt32 {
		recommend = math.MinInt32
	}
	modifiedNsec := int64(0)
	if !article.Modified().IsZero() {
		modifiedNsec = article.Modified().UnixNano()
	}
	return &apipb.Post{
		Index:         uint32(index + 1),
		Filename:      article.Filename(),
		RawDate:       article.Date(),
		NumRecommends: int32(recommend),
		Owner:         article.Owner(),
		Title:         article.Title(),
		ModifiedNsec:  modifiedNsec,
	}
}

func minInt(value, max int) int {
	if value > max {
		return max
	}
	return value
}
