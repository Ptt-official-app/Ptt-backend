package repository

import (
	"context"
	"fmt"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
)

var (
	logger = logging.NewLogger()
)

// Repository directly interacts with database via db handler.
type Repository interface {

	// board.go
	// GetBoards return all board record
	GetBoards(ctx context.Context) []bbs.BoardRecord
	// GetBoardArticle returns an article file in a specified board and filename
	GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error)
	// GetBoardArticleRecords returns article records of a board
	GetBoardArticleRecords(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error)
	// GetBoardTreasureRecords returns treasure article records of a board
	GetBoardTreasureRecords(ctx context.Context, boardID string, treasureIDs []string) ([]bbs.ArticleRecord, error)

	// user.go
	// GetUsers returns all user reords
	GetUsers(ctx context.Context) ([]bbs.UserRecord, error)
	// GetUserFavoriteRecords returns favorite records of a user
	GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error)
}

type repository struct {
	db           *bbs.DB
	userRecords  []bbs.UserRecord
	boardRecords []bbs.BoardRecord
}

func NewRepository(db *bbs.DB) (Repository, error) {
	userRecords, err := loadUserRecords(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load user records: %w", err)
	}

	boardRecords, err := loadBoardFile(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load board file: %w", err)
	}

	return &repository{
		db:           db,
		boardRecords: boardRecords,
		userRecords:  userRecords,
	}, nil
}
