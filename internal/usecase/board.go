package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/PichuChen/go-bbs"
)

type ArticleSearchCond struct {
	Title                           string
	Author                          string
	RecommendCountValue             int
	RecommendCountLessThan          int
	RecommendCountLessEqual         int
	RecommendCountEqual             int
	RecommendCountNotEqual          int
	RecommendCountGreaterThan       int
	RecommendCountGreaterEqual      int
	RecommendCountGreaterEqualIsSet bool
	RecommendCountLessEqualIsSet    bool
}

func (usecase *usecase) GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error) {
	for _, it := range usecase.repo.GetBoards(ctx) {
		if boardID == it.BoardId() {
			return it, nil
		}
	}
	return nil, fmt.Errorf("board record not found")
}

func (usecase *usecase) GetBoards(ctx context.Context, userID string) []bbs.BoardRecord {
	boards := make([]bbs.BoardRecord, 0)
	for _, board := range usecase.repo.GetBoards(ctx) {
		// TODO: Show Board by user level
		if board.IsClass() {
			continue
		}
		if !usecase.shouldShowOnUserLevel(board, userID) {
			continue
		}
		boards = append(boards, board)
	}
	return boards
}

func (usecase *usecase) GetPopularBoards(ctx context.Context) ([]bbs.BoardRecord, error) {
	boards := usecase.repo.GetBoards(ctx)
	// TODO: Sort boards in descending order by number_of_user

	// sort.Slice(boards, func(i, j int) bool {
	// 	return boards[i].UserNum > boards[j].UserNum
	// })
	if len(boards) < 100 {
		return boards, nil
	}
	return boards[:100], nil
}

func (usecase *usecase) GetClasses(ctx context.Context, userID, classID string) []bbs.BoardRecord {
	boards := make([]bbs.BoardRecord, 0)
	for _, board := range usecase.repo.GetBoards(ctx) {
		// TODO: Show Board by user level
		if !usecase.shouldShowOnUserLevel(board, userID) {
			continue
		}
		if board.ClassId() != classID {
			continue
		}
		// m := marshalBoardHeader(board)
		// if board.IsClass() {
		// 	m["id"] = fmt.Sprintf("%v", bid+1)
		// }
		boards = append(boards, board)
	}
	return boards
}

func (usecase *usecase) GetBoardArticles(ctx context.Context, boardID string, cond *ArticleSearchCond) []interface{} {
	var articles []bbs.ArticleRecord
	articleRecords, err := usecase.repo.GetBoardArticleRecords(ctx, boardID)
	if err != nil {
		usecase.logger.Warningf("open directory file error: %v", err)
		// The board may not contain any article
	}

	if len(strings.TrimSpace(cond.Title)) > 0 ||
		len(strings.TrimSpace(cond.Author)) > 0 ||
		cond.RecommendCountGreaterEqualIsSet ||
		cond.RecommendCountLessEqualIsSet {
		articles = searchArticles(articleRecords, cond)
	} else {
		articles = articleRecords
	}

	items := []interface{}{}
	for _, f := range articles {
		m := map[string]interface{}{
			"filename": f.Filename(),
			// Bug(pichu): f.Modified time will be 0 when file is vote
			"modified_time":   f.Modified(),
			"recommend_count": f.Recommend(),
			"post_date":       f.Date(),
			"title":           f.Title(),
			"money":           fmt.Sprintf("%v", f.Money()),
			"owner":           f.Owner(),
			// "aid": ""
			"url": getArticleURL(boardID, f.Filename()),
		}
		items = append(items, m)
	}
	return items
}

func (usecase *usecase) GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error) {
	buf, err := usecase.repo.GetBoardArticle(ctx, boardID, filename)
	if err != nil {
		return nil, fmt.Errorf("read file %s/%s error: %w", boardID, filename, err)
	}
	return buf, nil
}

func (usecase *usecase) GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{} {
	fileHeaders, err := usecase.repo.GetBoardTreasureRecords(ctx, boardID, treasuresID)
	if err != nil {
		usecase.logger.Warningf("open directory file error: %v", err)
		// The board may not contain any article
	}

	items := []interface{}{}
	for _, f := range fileHeaders {
		m := map[string]interface{}{
			"filename": f.Filename(),
			// Bug(pichu): f.Modified time will be 0 when file is vote
			"modified_time":   f.Modified(),
			"recommend_count": f.Recommend(),
			"post_date":       f.Date(),
			"title":           f.Title(),
			"money":           fmt.Sprintf("%v", f.Money()),
			"owner":           f.Owner(),
			// "aid": ""
			"url": getArticleURL(boardID, f.Filename()),
		}
		items = append(items, m)
	}
	return items
}

func (usecase *usecase) shouldShowOnUserLevel(board bbs.BoardRecord, userID string) bool {
	// TODO: Get user Level
	return true
}

func getArticleURL(boardId string, filename string) string {
	return fmt.Sprintf("https://ptt-app-dev-codingman.pichuchen.tw/bbs/%s/%s.html", boardId, filename)
}

func searchArticles(fileHeaders []bbs.ArticleRecord, cond *ArticleSearchCond) []bbs.ArticleRecord {
	var targetArticles []bbs.ArticleRecord

	for _, f := range fileHeaders {
		if !strings.Contains(strings.ToLower(f.Title()), strings.ToLower(cond.Title)) {
			continue
		}

		if !strings.Contains(strings.ToLower(f.Owner()), strings.ToLower(cond.Author)) {
			continue
		}

		if cond.RecommendCountGreaterEqualIsSet && f.Recommend() < cond.RecommendCountGreaterEqual {
			continue
		}

		if cond.RecommendCountLessEqualIsSet && f.Recommend() > cond.RecommendCountLessEqual {
			continue
		}

		targetArticles = append(targetArticles, f)
	}
	return targetArticles
}
