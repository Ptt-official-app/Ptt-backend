package repository

import (
	"context"
	"fmt"

	"github.com/PichuChen/go-bbs"
)

func (repo *repository) GetUsers(_ context.Context) []bbs.UserRecord {
	return repo.userRecords
}

func (repo *repository) GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error) {
	return repo.db.ReadUserFavoriteRecords(userID)
}

func loadUserRecords(db *bbs.DB) ([]bbs.UserRecord, error) {
	userRecords, err := db.ReadUserRecords()
	if err != nil {
		logger.Errorf("get user rec error: %v", err)
		return nil, fmt.Errorf("failed to read user records: %w", err)
	}
	return userRecords, nil
}
