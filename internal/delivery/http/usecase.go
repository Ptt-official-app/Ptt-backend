package http

import (
	"context"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	business "github.com/Ptt-official-app/Ptt-backend/internal/usecase"
	"github.com/Ptt-official-app/go-bbs"
)

// Usecase describes the business operations consumed by the HTTP delivery
// layer. It intentionally lives at the consumer instead of the usecase
// implementation package.
type Usecase interface {
	GetUserByID(ctx context.Context, userID string) (bbs.UserRecord, error)
	GetUserFavorites(ctx context.Context, userID string) ([]interface{}, error)
	AddUserFavorite(ctx context.Context, token, userID, favoriteType, boardID, title string) ([]interface{}, error)
	GetUserInformation(ctx context.Context, userID string) (map[string]interface{}, error)
	GetUserArticles(ctx context.Context, userID string) ([]interface{}, error)
	GetUserPreferences(ctx context.Context, userID string) (map[string]string, error)
	GetUserComments(ctx context.Context, userID string) ([]bbs.UserCommentRecord, error)
	GetUserDrafts(ctx context.Context, userID, draftID string) (bbs.UserDraft, error)
	UpdateUserDraft(ctx context.Context, userID, draftID string, text []byte) (bbs.UserDraft, error)
	DeleteUserDraft(ctx context.Context, userID, draftID string) error

	GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error)
	GetBoards(ctx context.Context, userID string) []bbs.BoardRecord
	CreateBoard(ctx context.Context, boardID, title string) (bbs.BoardRecord, error)
	GetPopularBoards(ctx context.Context) ([]bbs.BoardRecord, error)
	GetBoardPostsLimitation(ctx context.Context, boardID string) (*business.BoardPostLimitation, error)
	GetClasses(ctx context.Context, userID, classID string) ([]bbs.BoardRecord, error)
	GetBoardArticles(ctx context.Context, boardID string, offset, length uint, cond *business.ArticleSearchCond) []bbs.ArticleRecord
	GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error)
	GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{}
	CreateArticle(ctx context.Context, userID, boardID, title, article string) (bbs.ArticleRecord, error)
	DeleteArticle(ctx context.Context, token, boardID, filename string) error
	GetRawArticle(boardID, filename string) (string, error)

	CreateAccessTokenWithUsername(username string) string
	GetUserIDFromToken(token string) (string, error)
	RecordLogin(userID, ip string)
	CheckPermission(token string, permissionID []business.Permission, userInfo map[string]string) error

	GetPopularArticles(ctx context.Context) ([]repository.PopularArticleRecord, error)
	AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (repository.PushRecord, error)
	ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (repository.ForwardArticleToBoardRecord, error)
	ForwardArticleToEmail(ctx context.Context, userID, boardID, filename, email string) error
	UpdateUsefulness(ctx context.Context, userID, boardID, filename, appendType string) (repository.PushRecord, error)
}

type articleURLProvider interface {
	GetArticleURL(boardID, filename string) string
}
