package repository

import (
	"context"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
)

var (
	logger = logging.NewLogger()
)

type BoardRepository interface {
	GetBoards(ctx context.Context) []bbs.BoardRecord
	GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error)
	GetBoardArticleRecords(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error)
	GetBoardTreasureRecords(ctx context.Context, boardID string, treasureIDs []string) ([]bbs.ArticleRecord, error)
}

type UserRepository interface {
	GetUsers(ctx context.Context) []bbs.UserRecord
	GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error)
}
