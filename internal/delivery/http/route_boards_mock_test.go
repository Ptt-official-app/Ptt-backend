package http

import (
	"context"
	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (usecase *MockUsecase) GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error) {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoards(ctx context.Context, userID string) []bbs.BoardRecord {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoardArticles(ctx context.Context, boardID string, cond *usecase.ArticleSearchCond) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"filename":        "test-001",
			"modified_time":   "2009-01-01T12:59:59Z",
			"recommend_count": 9,
			"post_date":       "2009-01-01",
			"title":           "post for testing",
			"money":           "10",
			"owner":           "tester",
			"url":             "http://test/test-001.html",
		},
	}
}

func (usecase *MockUsecase) GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error) {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{} {
	panic("Not implemented")
}
