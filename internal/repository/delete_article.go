package repository

import (
	"context"
	"fmt"
)

func (repo *repository) DeleteBoardArticle(_ context.Context, boardID, filename, deletedBy string) error {
	if err := repo.db.DeleteBoardArticle(boardID, filename, deletedBy); err != nil {
		return fmt.Errorf("delete board article %s/%s: %w", boardID, filename, err)
	}
	repo.popularArticles.invalidate()
	return nil
}
