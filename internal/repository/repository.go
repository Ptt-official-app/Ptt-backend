package repository

import (
	"context"
	"fmt"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/db"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
)

var (
	logger = logging.NewLogger()
)

type Repository interface {

	// board.go
	GetBoards(ctx context.Context) []bbs.BoardRecord
	GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error)
	GetBoardArticleRecords(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error)
	GetBoardTreasureRecords(ctx context.Context, boardID string, treasureIDs []string) ([]bbs.ArticleRecord, error)

	// user.go
	GetUsers(ctx context.Context) []bbs.UserRecord
	GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error)
}

type repository struct {
	db           db.DB
	userRecords  []bbs.UserRecord
	boardRecords []bbs.BoardRecord
}

func NewRepository(db db.DB) (Repository, error) {
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
