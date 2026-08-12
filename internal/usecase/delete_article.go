package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/go-bbs"
)

var (
	ErrArticleNotFound        = errors.New("article not found")
	ErrDeleteArticleForbidden = errors.New("article deletion forbidden")
)

func (usecase *usecase) DeleteArticle(ctx context.Context, token, boardID, filename string) error {
	userID, err := usecase.GetUserIDFromToken(token)
	if err != nil {
		return fmt.Errorf("authenticate article deletion: %w", err)
	}

	board, err := usecase.GetBoardByID(ctx, boardID)
	if err != nil {
		return fmt.Errorf("get board %s: %w", boardID, err)
	}

	article, err := usecase.getArticleRecord(ctx, boardID, filename)
	if err != nil {
		return err
	}

	user, err := usecase.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user %s: %w", userID, err)
	}

	if !canDeleteArticle(user, board, article) {
		return fmt.Errorf("%w: user %s cannot delete %s/%s", ErrDeleteArticleForbidden, userID, boardID, filename)
	}

	if err := usecase.repo.DeleteBoardArticle(ctx, boardID, filename, userID); err != nil {
		return fmt.Errorf("delete article: %w", err)
	}
	return nil
}

func (usecase *usecase) getArticleRecord(ctx context.Context, boardID, filename string) (bbs.ArticleRecord, error) {
	records, err := usecase.repo.GetBoardArticleRecords(ctx, boardID, 0, ^uint(0))
	if err != nil {
		return nil, fmt.Errorf("list board articles: %w", err)
	}
	for _, article := range records {
		if article != nil && article.Filename() == filename {
			return article, nil
		}
	}
	return nil, fmt.Errorf("%w: %s/%s", ErrArticleNotFound, boardID, filename)
}

func canDeleteArticle(user bbs.UserRecord, board bbs.BoardRecord, article bbs.ArticleRecord) bool {
	if user == nil || board == nil || article == nil {
		return false
	}
	if repository.UserIsSYSOP(user) {
		return true
	}

	for _, bm := range board.BM() {
		if strings.EqualFold(bm, user.UserID()) {
			return true
		}
	}

	if !strings.EqualFold(article.Owner(), user.UserID()) {
		return false
	}
	if settings, ok := board.(bbs.BoardRecordSettings); ok && settings.IsNoSelfDeletePost() {
		return false
	}
	return true
}
