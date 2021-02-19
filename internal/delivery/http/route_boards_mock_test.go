package http

import (
	"context"
	"github.com/Ptt-official-app/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (usecase *MockUsecase) GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error) {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoards(ctx context.Context, userID string) []bbs.BoardRecord {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoardArticles(ctx context.Context, boardID string, cond *usecase.ArticleSearchCond) []interface{} {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error) {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{} {
	panic("Not implemented")
}
