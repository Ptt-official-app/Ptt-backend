package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/Ptt-official-app/go-bbs"
)

// Repository directly interacts with database via db handler.
type Repository interface {

	// board.go
	// GetBoards return all board record
	GetBoards(ctx context.Context) []bbs.BoardRecord
	// CreateBoard creates a new board and updates the in-process board cache.
	CreateBoard(ctx context.Context, boardID, title string) (bbs.BoardRecord, error)
	// GetBoardArticle returns an article file in a specified board and filename
	GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error)
	// GetBoardArticleRecords returns article records of a board
	GetBoardArticleRecords(ctx context.Context, boardID string, offset, length uint) ([]bbs.ArticleRecord, error)
	// GetBoardTreasureRecords returns treasure article records of a board
	GetBoardTreasureRecords(ctx context.Context, boardID string, treasureIDs []string) ([]bbs.ArticleRecord, error)

	// user.go
	GetUsers(ctx context.Context) ([]bbs.UserRecord, error)
	GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error)
	GetUserArticles(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error)
	GetUserPreferences(ctx context.Context, userID string) (map[string]string, error)
	GetUserComments(ctx context.Context, userID string) ([]bbs.UserCommentRecord, error)
	GetUserDrafts(ctx context.Context, userID, draftID string) (bbs.UserDraft, error)
	UpdateUserDraft(ctx context.Context, userID, draftID string, text []byte) (bbs.UserDraft, error)
	DeleteUserDraft(ctx context.Context, userID, draftID string) error

	// article.go
	GetPopularArticles(ctx context.Context) ([]PopularArticleRecord, error)
	AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (PushRecord, error)
	CreateArticle(ctx context.Context, userID, boardID, title, article string) (bbs.ArticleRecord, error)
	GetRawArticle(boardID, filename string) (string, error)
	ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (ForwardArticleToBoardRecord, error)
}

type repository struct {
	db              *bbs.DB
	userRecords     []bbs.UserRecord
	boardMu         sync.RWMutex
	boardRecords    []bbs.BoardRecord
	popularArticles popularArticlesCache
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
