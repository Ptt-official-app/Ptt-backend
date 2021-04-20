package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ptt-official-app/go-bbs"
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

type BoardPostLimitation struct {
	PostsLimit   uint8
	LoginsLimit  uint8
	BadPostLimit uint8
}

func (usecase *usecase) GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error) {
	for _, it := range usecase.repo.GetBoards(ctx) {
		if boardID == it.BoardID() {
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
	// Use GetBoards to obtain data and use conditions to filter
	// TODO:GetBoards need add return error
	allBoards := usecase.repo.GetBoards(ctx)
	if len(allBoards) == 0 {
		usecase.logger.Warningf("GetBoards Error")
	}

	filtedBoards := shouldBeDisplayOnPouplarList(&allBoards)

	// TODO:Add condition to sort
	// sort.Slice(filtedBoards, func(i, j int) bool {
	// 	return (*filtedBoards)[i].UserEntered > (*filtedBoards)[j].UserEntered
	// })

	if len(*filtedBoards) < 100 {
		return *filtedBoards, nil
	}

	return (*filtedBoards)[:100], nil
}

func shouldBeDisplayOnPouplarList(allBoards *[]bbs.BoardRecord) (filtedBoards *[]bbs.BoardRecord) {
	for index := range *allBoards {
		// Initially filter boards by board status or other values
		// TODO:Need to add filter conditions,here is an example
		if (*allBoards)[index].ClassID() != "" {
			*filtedBoards = append(*filtedBoards, (*allBoards)[index])
		}
		// TODO:Add other conditions
	}
	return
}

func (usecase *usecase) GetBoardPostsLimitation(ctx context.Context, boardID string) (*BoardPostLimitation, error) {
	postLimit, err := usecase.repo.GetBoardPostsLimit(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("get board %s posts limit error: %w", boardID, err)
	}

	loginsLimit, err := usecase.repo.GetBoardLoginsLimit(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("get board %s logins limit error: %w", boardID, err)
	}

	badPostLimit, err := usecase.repo.GetBoardBadPostLimit(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("get board %s bad posts limit error: %w", boardID, err)
	}

	return &BoardPostLimitation{
		PostsLimit:   postLimit.PostLimitPosts(),
		LoginsLimit:  loginsLimit.PostLimitLogins(),
		BadPostLimit: badPostLimit.PostLimitBadPost(),
	}, nil
}

func (usecase *usecase) GetClasses(ctx context.Context, userID, classID string) []bbs.BoardRecord {
	boards := make([]bbs.BoardRecord, 0)
	for _, board := range usecase.repo.GetBoards(ctx) {
		// TODO: Show Board by user level
		if !usecase.shouldShowOnUserLevel(board, userID) {
			continue
		}
		if board.ClassID() != classID {
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

func getArticleURL(boardID string, filename string) string {
	// TODO: generate article url by config file
	return fmt.Sprintf("https://pttapp.cc/bbs/%s/%s.html", boardID, filename)
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
