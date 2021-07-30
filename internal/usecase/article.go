package usecase

import (
	"context"
	"fmt"
	"time"

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

func (usecase *usecase) UpdateUsefulness(ctx context.Context, userID, boardID, filename, appendType string) (repository.PushRecord, error) {
	return nil, nil
}

// AppendComment append comment to specific article
func (usecase *usecase) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (repository.PushRecord, error) {
	result, err := usecase.repo.AppendComment(ctx, userID, boardID, filename, appendType, text)

	return result, err
}

// ForwardArticleToBoard returns forwarding to board results
func (usecase *usecase) ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (repository.ForwardArticleToBoardRecord, error) {
	token := usecase.CreateAccessTokenWithUsername(userID)

	err := usecase.checkForwardArticlePermission(token,
		map[string]string{
			"board_id":   boardID,
			"article_id": filename,
		})
	if err != nil {
		return nil, fmt.Errorf("checkForwardArticlePermission error: %w", err)
	}

	forwardArticle, err := usecase.repo.ForwardArticleToBoard(ctx, userID, boardID, filename, boardName)
	if err != nil {
		return nil, fmt.Errorf("ForwardArticleToBoard error: %w", err)
	}

	return forwardArticle, err
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
	return usecase.mailProvider.Send(email, title, userID, buffer)
}

// CreateArticle create a new article on a board
func (usecase *usecase) CreateArticle(ctx context.Context, userID, boardID, title, article string) (bbs.ArticleRecord, error) {
	now := time.Now()
	bbsDate := fmt.Sprintf("%02d/%02d", now.Month(), now.Day())

	err := usecase.repo.CreateArticle(ctx, userID, boardID, title, article)
	if err != nil {
		return nil, err
	}

	articles, err := usecase.repo.GetBoardArticleRecords(ctx, boardID)

	if err != nil {
		return nil, fmt.Errorf("get article records failed: %w", err)
	}

	for i := 0; i < len(articles); i++ {
		if articles[i] != nil && articles[i].Owner() == userID && articles[i].Title() == title && articles[i].Date() == bbsDate {
			return articles[i], nil
		}
	}

	return nil, fmt.Errorf("get author latest article records failed: %w", err)
}

func (usecase *usecase) GetArticleURL(boardID string, filename string) string {
	// TODO: generate article url by config file
	return fmt.Sprintf("https://pttapp.cc/bbs/%s/%s.html", boardID, filename)
}
