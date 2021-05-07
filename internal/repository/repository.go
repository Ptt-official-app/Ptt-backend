package repository

import (
	"context"
	"fmt"

	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/go-bbs"
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
	// GetBoardPostsLimit returns posts limited record of a board
	// TODO: replace PostsLimitedBoardRecord with real bbs record
	GetBoardPostsLimit(ctx context.Context, boardID string) (PostsLimitedBoardRecord, error)
	// GetBoardLoginsLimit returns logins limited record of a board
	// TODO: replace LoginsLimitedBoardRecord with real bbs record
	GetBoardLoginsLimit(ctx context.Context, boardID string) (LoginsLimitedBoardRecord, error)
	// GetBoardBadPostLimit returns bad posts limited record of a board
	// TODO: replace BadPostLimitedBoardRecord with real bbs record
	GetBoardBadPostLimit(ctx context.Context, boardID string) (BadPostLimitedBoardRecord, error)

	// user.go
	// GetUsers returns all user records
	GetUsers(ctx context.Context) ([]bbs.UserRecord, error)
	// GetUserFavoriteRecords returns favorite records of a user
	GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error)
	// GetUserArticles returns user's articles
	GetUserArticles(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error)
	// GetUserPreferences returns user's preferences
	// TODO: replace UserPreferencesRecord with real bbs record
	GetUserPreferences(ctx context.Context, userID string) (map[string]string, error)
	// GetUserComments return user's history comments
	// TODO: return a slice of concrete type not interface
	GetUserComments(ctx context.Context, userID string) ([]interface{}, error)
	// GetUserDrafts returns user's draft according to draftID
	GetUserDrafts(ctx context.Context, userID, draftID string) (UserDraft, error)
	// UpdateUserDraft updates user's draft according to draftID
	UpdateUserDraft(ctx context.Context, userID, draftID string, text []byte) (UserDraft, error)
	// DeleteUserDraft deletes user's draft according to draftID
	DeleteUserDraft(ctx context.Context, userID, draftID string) error

	// article.go
	// GetPopularArticles returns all popular articles
	GetPopularArticles(ctx context.Context) ([]PopularArticleRecord, error)
	// AppendComment returns comment details
	AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (PushRecord, error)
	// CreateArticle
	// TODO: return result from bbs response
	CreateArticle(ctx context.Context, userID, boardID, title, article string) error
	// ForwardArticleToBoard returns forwarding to board results
	ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (ForwardArticleToBoardRecord, error)
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
