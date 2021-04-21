package usecase

import (
	"context"
	"fmt"

	"github.com/Ptt-official-app/go-bbs"

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
	//result, err := usecase.repo.AppendComment(ctx, userID, boardID, filename, appendType, text)

	//return result, err
	return nil, nil
}

// ForwardArticleToBoard returns forwarding to board results
func (usecase *usecase) ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// ForwardArticleToEmail returns forwarding to email results
func (usecase *usecase) ForwardArticleToEmail(ctx context.Context, userID, boardID, filename, email string) error {
	articleRecords, err := usecase.repo.GetBoardArticleRecords(ctx, boardID)
	if err != nil {
		return fmt.Errorf("GetBoardArticleRecords error: %w", err)
	}
	var title string
	for _, article := range articleRecords {
		if article.Filename() == filename {
			title = article.Title()
			break
		}
	}
	if title == "" {
		return fmt.Errorf("cannot find article %s", filename)
	}
	buffer, err := usecase.repo.GetBoardArticle(ctx, boardID, filename)
	if err != nil {
		return fmt.Errorf("GetBoardArticle error: %w", err)
	}
	return usecase.mail.Send(email, title, userID, buffer)
}

// CreateArticle create a new article on a board
func (usecase *usecase) CreateArticle(ctx context.Context, userID, boardID, title, article string) (bbs.ArticleRecord, error) {
	return nil, nil
}
