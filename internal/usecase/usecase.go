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

type BoardUsecase interface {
	GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error)
	GetBoards(ctx context.Context, userID string) []bbs.BoardRecord
	GetClasses(ctx context.Context, userID, classID string) []bbs.BoardRecord
	// FIXME: use concrete type rather than []interface{}
	GetBoardArticles(ctx context.Context, boardID string) []interface{}
	GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error)
	// FIXME: use concrete type rather than []interface{}
	GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{}
}
