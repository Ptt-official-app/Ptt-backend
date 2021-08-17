package usecase

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

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
	cond := &ArticleSearchCond{
		Filename: filename,
	}
	articleRecord := usecase.GetBoardArticles(ctx, boardID, cond)

	if articleRecord[0].Owner() == userID {
		return nil, fmt.Errorf("Owners cannot push their own article")
	}

	article, err := usecase.GetBoardArticle(ctx, boardID, filename)
	if err != nil {
		return nil, fmt.Errorf("GetBoardArticle error: %w", err)
	}

	articleStr := string(article)

	cur := 0
	numRecommend := 0
	for {
		cur = strings.Index(articleStr[cur:], userID)
		if cur < 0 {
			break
		}

		r, _ := utf8.DecodeLastRuneInString(articleStr[:cur])
		if r == utf8.RuneError {
			return nil, fmt.Errorf("DecodeLastRuneError in usecase UpdateUsefulness")
		}
		if numRecommend < 1 && string(r) == "\u2191" {
			numRecommend++
		}

		if numRecommend > -1 && string(r) == "\u2193" {
			numRecommend--
		}
	}

	if (appendType == "\u2191" && numRecommend == 1) || (appendType == "\u2193" && numRecommend == -1) {
		return nil, fmt.Errorf("Cannot push this time")
	} else {
		p := usecase.repo.GetPushRecord(ctx, userID, boardID, filename, appendType)
	}

	return p, nil
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
	err := usecase.repo.CreateArticle(ctx, userID, boardID, title, article)
	if err != nil {
		return nil, err
	}

	articles, err := usecase.repo.GetBoardArticleRecords(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("get article records failed: %w", err)
	}

	return articles[0], nil // todo: get first for temporary
}

func (usecase *usecase) GetArticleURL(boardID string, filename string) string {
	// TODO: generate article url by config file
	return fmt.Sprintf("https://pttapp.cc/bbs/%s/%s.html", boardID, filename)
}
