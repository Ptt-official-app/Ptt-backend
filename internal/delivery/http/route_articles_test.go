package http

import (
	"context"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (usecase *MockUsecase) GetPopularArticles(ctx context.Context) ([]repository.PopularArticleRecord, error) {
	return []repository.PopularArticleRecord{}, nil
}

func (usecase *MockUsecase) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (usecase *MockUsecase) ForwardArticle(ctx context.Context, userID, boardID, filename string, to usecase.Forward) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
