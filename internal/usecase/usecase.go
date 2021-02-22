package usecase

import (
	"context"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

// Usecase is the implementation of backend business logic.
type Usecase interface {
	// user.go
	// GetUserByID returns a user of userID
	GetUserByID(ctx context.Context, userID string) (bbs.UserRecord, error)
	// GetUserFavorites returns favorite records of a user
	GetUserFavorites(ctx context.Context, userID string) ([]interface{}, error) // FIXME: use concrete type rather than []interface{}
	// GetUserInformation returns user info of a user
	GetUserInformation(ctx context.Context, userID string) (map[string]interface{}, error) // FIXME: use concrete type rather than map[string]interface{}

	// board.go
	// GetBoardByID returns board record of board id
	GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error)
	// GetBoards returns all board records
	GetBoards(ctx context.Context, userID string) []bbs.BoardRecord
	// GetPopularBoards returns top 100 popular board records
	GetPopularBoards(ctx context.Context) ([]bbs.BoardRecord, error)
	// GetClasses returns board records in a class
	GetClasses(ctx context.Context, userID, classID string) []bbs.BoardRecord
	// GetBoardArticles returns articles of a board
	GetBoardArticles(ctx context.Context, boardID string, cond *ArticleSearchCond) []interface{} // FIXME: use concrete type rather than []interface{}
	// GetBoardArticle returns an article file given board id and file name
	GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error)
	// GetBoardTreasures returns treasures of a board
	GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{} // FIXME: use concrete type rather than []interface{}

	// token.go
	// CreateAccessTokenWithUsername creates access token for a user
	CreateAccessTokenWithUsername(username string) string
	// GetUserIdFromToken retrieves user id by token
	GetUserIdFromToken(token string) (string, error)
	// CheckPermission checks permissions
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
