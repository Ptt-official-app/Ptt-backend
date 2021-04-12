package usecase

import (
	"context"
	"fmt"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

// GetPopularArticles returns articles by descending comment_count
func (usecase *usecase) GetPopularArticles(ctx context.Context) ([]repository.PopularArticleRecord, error) {
	articles, err := usecase.repo.GetPopularArticles(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetPopularArticles error: %w", err)
	}
	return articles, nil
}

// AppendComment append comment to specific article
func (usecase *usecase) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// ForwardArticleToBoard returns forwarding to board results
func (usecase *usecase) ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, toBoard string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// ForwardArticleToEmail returns forwarding to email results
func (usecase *usecase) ForwardArticleToEmail(ctx context.Context, userID, boardID, filename, toEmail string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
