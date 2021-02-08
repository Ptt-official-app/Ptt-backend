package repository

import (
	"context"
	"testing"

	"github.com/PichuChen/go-bbs"
)

type MockRepository struct{}

// repository/board.go
func (repo *MockRepository) GetBoards(ctx context.Context) []bbs.BoardRecord {
	return []bbs.BoardRecord{}
}
func (repo *MockRepository) GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error) {
	return []byte{}, nil
}
func (repo *MockRepository) GetBoardArticleRecords(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error) {
	return []bbs.ArticleRecord{}, nil
}

func (repo *MockRepository) GetBoardTreasureRecords(ctx context.Context, boardID string, treasureIDs []string) ([]bbs.ArticleRecord, error) {
	return []bbs.ArticleRecord{}, nil
}

func TestGetUsers(t *testing.T) {

	repo := repository{
		userRecords:  []bbs.UserRecord{},
		boardRecords: []bbs.BoardRecord{},
	}

	actual := repo.GetUsers(context.TODO())
	if actual == nil {
		t.Errorf("GetUsers got %v, expected not equal nil", actual)
	}

}
