package http

import (
	"context"

	"github.com/Ptt-official-app/go-bbs"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

func (usecase *MockUsecase) GetPopularArticles(ctx context.Context) ([]repository.PopularArticleRecord, error) {
	return []repository.PopularArticleRecord{}, nil
}

func (usecase *MockUsecase) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (repository.PushRecord, error) {
	return nil, nil
}

func (usecase *MockUsecase) ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (repository.ForwardArticleToBoardRecord, error) {
	return nil, nil
}

func (usecase *MockUsecase) ForwardArticleToEmail(ctx context.Context, userID, boardID, filename, email string) error {
	return nil
}

func (usecase *MockUsecase) CreateArticle(ctx context.Context, userID, boardID, title, article string) (bbs.ArticleRecord, error) {
	return nil, nil
}
