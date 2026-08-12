package repository

import (
	"context"
	"fmt"

	"github.com/Ptt-official-app/go-bbs"
)

func (repo *repository) AddUserFavorite(_ context.Context, userID string, options bbs.FavoriteCreateOptions) (bbs.FavoriteRecord, error) {
	record, err := repo.db.AddUserFavorite(userID, options)
	if err != nil {
		return nil, fmt.Errorf("add favorite for user %s: %w", userID, err)
	}
	return record, nil
}
