package http

import (
	"context"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

func (usecase *MockUsecase) GetPopularArticles(ctx context.Context) ([]repository.PopularArticleRecord, error) {
	return []repository.PopularArticleRecord{}, nil
}

func (usecase *MockUsecase) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (usecase *MockUsecase) ForwardArticle(ctx context.Context, userID, boardID, filename, toEmail, toBoard string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
