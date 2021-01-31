package mock

import (
	"context"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

// implements repository.Repository
type MockRepository struct {
}

func NewMockRepository() repository.Repository {
	return &MockRepository{}
}

// repository/board.go
func (repo *MockRepository) GetBoards(ctx context.Context) []bbs.BoardRecord {
	panic("Not implemented")
}
func (repo *MockRepository) GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error) {
	panic("Not implemented")
}
func (repo *MockRepository) GetBoardArticleRecords(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error) {
	panic("Not implemented")
}

func (repo *MockRepository) GetBoardTreasureRecords(ctx context.Context, boardID string, treasureIDs []string) ([]bbs.ArticleRecord, error) {
	panic("Not implemented")
}

// repository/user.go
func (repo *MockRepository) GetUsers(ctx context.Context) []bbs.UserRecord {
	panic("Not implemented")
}
func (repo *MockRepository) GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error) {
	panic("Not implemented")
}
