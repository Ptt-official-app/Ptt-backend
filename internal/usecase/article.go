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

type Forward interface {
	Forward() (map[string]interface{}, error)
}
type ForwardToBoard struct {
	Board string
}

func (b *ForwardToBoard) Forward() (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

type ForwardToEmail struct {
	Email string
}

func (b *ForwardToEmail) Forward() (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// ForwardArticle returns forwarding to board results
func (usecase *usecase) ForwardArticle(ctx context.Context, userID, boardID, filename string, to Forward) (map[string]interface{}, error) {
	return to.Forward()
}
