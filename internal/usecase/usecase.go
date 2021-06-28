package usecase

import (
	"context"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/mail"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/go-bbs"
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
	// GetUserArticles returns user's articles
	GetUserArticles(ctx context.Context, userID string) ([]interface{}, error) // FIXME: use concrete type rather than []interface{}
	GetUserPreferences(ctx context.Context, userID string) (map[string]string, error)
	// GetUserComments returns history comments of a user
	GetUserComments(ctx context.Context, userID string) ([]bbs.UserCommentRecord, error) // FIXME: use concrete type from go-bbs instead of []interface{}
	// GetUserDrafts returns user's draft by given draft id
	GetUserDrafts(ctx context.Context, userID, draftID string) (bbs.UserDraft, error)
	// UpdateUserDraft returns updated content
	UpdateUserDraft(ctx context.Context, userID, draftID string, text []byte) (bbs.UserDraft, error)
	DeleteUserDraft(ctx context.Context, userID, draftID string) error

	// board.go
	// GetBoardByID returns board record of board id
	GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error)
	// GetBoards returns all board records
	GetBoards(ctx context.Context, userID string) []bbs.BoardRecord
	// GetPopularBoards returns top 100 popular board records
	GetPopularBoards(ctx context.Context) ([]bbs.BoardRecord, error)
	// GetBoardPostsLimition returns all posts limit of a board
	GetBoardPostsLimitation(ctx context.Context, boardID string) (*BoardPostLimitation, error)
	// GetClasses returns board records in a class
	GetClasses(ctx context.Context, userID, classID string) ([]bbs.BoardRecord, error)
	// GetBoardArticles returns articles of a board
	GetBoardArticles(ctx context.Context, boardID string, cond *ArticleSearchCond) []interface{} // FIXME: use concrete type rather than []interface{}
	// GetBoardArticle returns an article file given board id and file name
	GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error)
	// GetBoardTreasures returns treasures of a board
	GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{} // FIXME: use concrete type rather than []interface{}
	// CreatePost create a new post
	CreateArticle(ctx context.Context, userID, boardID, title, article string) (bbs.ArticleRecord, error)

	// token.go
	// CreateAccessTokenWithUsername creates access token for a user
	CreateAccessTokenWithUsername(username string) string
	// GetUserIDFromToken retrieves user id by token
	GetUserIDFromToken(token string) (string, error)
	// CheckPermission checks permissions
	CheckPermission(token string, permissionID []Permission, userInfo map[string]string) error // FIXME: use concrete type rather than map[string]string

	// article.go
	// GetPopularArticles returns all popular articles
	GetPopularArticles(ctx context.Context) ([]repository.PopularArticleRecord, error)
	AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (repository.PushRecord, error)
	// ForwardArticleToBoard returns forwarding to board results
	ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (repository.ForwardArticleToBoardRecord, error)
	// ForwardArticleToEmail returns forwarding to email results
	ForwardArticleToEmail(ctx context.Context, userID, boardID, filename, email string) error
	// UpdateUsefulness update article usefulness
	UpdateUsefulness(ctx context.Context, userID, boardID, filename, appendType string) (repository.PushRecord, error)

	// mail.go
	UpdateMail(mail mail.Mail) error
}

type usecase struct {
	logger       logging.Logger
	globalConfig *config.Config
	repo         repository.Repository
	mailProvider mail.Mail
}

func NewUsecase(globalConfig *config.Config, repo repository.Repository) Usecase {
	mailProvider, _ := mail.NewMailProvider(globalConfig.MailDriver)
	return &usecase{
		logger:       logging.NewLogger(),
		globalConfig: globalConfig,
		repo:         repo,
		mailProvider: mailProvider,
	}
}
