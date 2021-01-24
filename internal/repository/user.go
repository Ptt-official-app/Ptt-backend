package repository

import (
	"context"
	"fmt"

	"github.com/PichuChen/go-bbs"
)

type userRepository struct {
	db          *bbs.DB
	userRecords []bbs.UserRecord
}

func NewUserRepository(db *bbs.DB) (UserRepository, error) {
	userRecords, err := loadUserRecords(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load user records: %w", err)
	}

	return &userRepository{
		db:          db,
		userRecords: userRecords,
	}, nil
}

func (u *userRepository) GetUsers(_ context.Context) []bbs.UserRecord {
	return u.userRecords
}

func (u *userRepository) GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error) {
	return u.db.ReadUserFavoriteRecords(userID)
}

func loadUserRecords(db *bbs.DB) ([]bbs.UserRecord, error) {
	userRecords, err := db.ReadUserRecords()
	if err != nil {
		logger.Errorf("get user rec error: %v", err)
		return nil, fmt.Errorf("failed to read user records: %w", err)
	}
	return userRecords, nil
}
