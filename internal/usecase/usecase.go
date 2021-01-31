package usecase

import (
	"context"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

type Usecase interface {
	// user.go
	GetUserByID(ctx context.Context, userID string) (bbs.UserRecord, error)
	GetUserFavorites(ctx context.Context, userID string) ([]interface{}, error)            // FIXME: use concrete type rather than []interface{}
	GetUserInformation(ctx context.Context, userID string) (map[string]interface{}, error) // FIXME: use concrete type rather than map[string]interface{}

	// board.go
	GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error)
	GetBoards(ctx context.Context, userID string) []bbs.BoardRecord
	GetClasses(ctx context.Context, userID, classID string) []bbs.BoardRecord
	GetBoardArticles(ctx context.Context, boardID string) []interface{} // FIXME: use concrete type rather than []interface{}
	GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error)
	GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{} // FIXME: use concrete type rather than []interface{}

	// token.go
	CreateAccessTokenWithUsername(username string) string
	GetUserIdFromToken(token string) (string, error)
	CheckPermission(token string, permissionId []Permission, userInfo map[string]string) error // FIXME: use concrete type rather than map[string]string
}

type usecase struct {
	logger       logging.Logger
	globalConfig *config.Config
	repo         repository.Repository
}

func NewUsecase(globalConfig *config.Config, repo repository.Repository) Usecase {
	return &usecase{
		logger:       logging.NewLogger(),
		globalConfig: globalConfig,
		repo:         repo,
	}
}
