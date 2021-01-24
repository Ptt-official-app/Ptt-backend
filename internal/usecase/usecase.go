package usecase

import (
	"context"

	"github.com/PichuChen/go-bbs"
)

type UserUsecase interface {
	GetUserByID(ctx context.Context, userID string) (bbs.UserRecord, error)
	// FIXME: use concrete type rather than []interface{}
	GetUserFavorites(ctx context.Context, userID string) ([]interface{}, error)
	// FIXME: use concrete type rather than map[string]interface{}
	GetUserInformation(ctx context.Context, userID string) (map[string]interface{}, error)
}
