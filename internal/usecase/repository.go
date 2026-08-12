package usecase

import (
	"context"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/go-bbs"
)

// Repository describes the persistence operations consumed by the business
// logic. The interface lives here, at the consumer, so repository
// implementations do not need to own or anticipate their consumers' contract.
type Repository interface {
	GetBoards(ctx context.Context) []bbs.BoardRecord
	CreateBoard(ctx context.Context, boardID, title string) (bbs.BoardRecord, error)
	GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error)
	DeleteBoardArticle(ctx context.Context, boardID, filename, deletedBy string) error
	GetBoardArticleRecords(ctx context.Context, boardID string, offset, length uint) ([]bbs.ArticleRecord, error)
	GetBoardTreasureRecords(ctx context.Context, boardID string, treasureIDs []string) ([]bbs.ArticleRecord, error)

	GetUsers(ctx context.Context) ([]bbs.UserRecord, error)
	GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error)
	AddUserFavorite(ctx context.Context, userID string, options bbs.FavoriteCreateOptions) (bbs.FavoriteRecord, error)
	GetUserArticles(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error)
	GetUserPreferences(ctx context.Context, userID string) (map[string]string, error)
	GetUserComments(ctx context.Context, userID string) ([]bbs.UserCommentRecord, error)
	GetUserDrafts(ctx context.Context, userID, draftID string) (bbs.UserDraft, error)
	UpdateUserDraft(ctx context.Context, userID, draftID string, text []byte) (bbs.UserDraft, error)
	DeleteUserDraft(ctx context.Context, userID, draftID string) error

	GetPopularArticles(ctx context.Context) ([]repository.PopularArticleRecord, error)
	AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (repository.PushRecord, error)
	CreateArticle(ctx context.Context, userID, boardID, title, article string) (bbs.ArticleRecord, error)
	GetRawArticle(boardID, filename string) (string, error)
	ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (repository.ForwardArticleToBoardRecord, error)
}
