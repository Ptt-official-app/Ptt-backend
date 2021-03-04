package http

import (
	"context"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

func (usecase *MockUsecase) GetPopularArticles(ctx context.Context) ([]repository.PopularArticleRecord, error) {
	return []repository.PopularArticleRecord{}, nil
}
