package usecase

import (
	"context"

	"github.com/Ptt-official-app/go-bbs"
)

func (repo *MockRepository) AddUserFavorite(ctx context.Context, userID string, options bbs.FavoriteCreateOptions) (bbs.FavoriteRecord, error) {
	return nil, nil
}
