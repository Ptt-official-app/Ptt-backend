package usecase

import (
	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"

	"github.com/PichuChen/go-bbs"

	"context"
	"testing"
	"time"
)

// implements repository.Repository
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

// repository/user.go
func (repo *MockRepository) GetUsers(ctx context.Context) []bbs.UserRecord {
	ret := []bbs.UserRecord{}
	ret = append(ret, &MockUser{
		userID: "pichu",
	})
	return ret
}
func (repo *MockRepository) GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error) {
	return []bbs.FavoriteRecord{}, nil
}

type MockUser struct {
	userID string
}

func (u *MockUser) UserId() string                       { return u.userID }
func (u *MockUser) HashedPassword() string               { return "" }
func (u *MockUser) VerifyPassword(password string) error { return nil }
func (u *MockUser) Nickname() string                     { return "" }
func (u *MockUser) RealName() string                     { return "" }
func (u *MockUser) NumLoginDays() int                    { return 0 }
func (u *MockUser) NumPosts() int                        { return 0 }
func (u *MockUser) Money() int                           { return 0 }
func (u *MockUser) LastLogin() time.Time                 { return time.Unix(0, 0) }
func (u *MockUser) LastHost() string                     { return "" }

func TestGetUserByID(t *testing.T) {

	resp := &MockRepository{}

	usecase := NewUsecase(&config.Config{}, resp)

	rec, err := usecase.GetUserByID(context.TODO(), "not-exist-user-id")
	if err == nil {
		t.Errorf("getUserByID with not-exist-user-id excepted not nil error, got nil")
		return
	}

	if rec != nil {
		t.Errorf("getUserByID with not-exist-user-id excepted nil, got %v", rec)
		return
	}

	rec, err = usecase.GetUserByID(context.TODO(), "pichu")
	if err != nil {
		t.Errorf("getUserByID with pichu excepted err == nil, got %v", err)
		return
	}

	if rec.UserId() != "pichu" {
		t.Errorf("getUserByID with pichu excepted userid: pichu, got %v", rec.UserId())
		return
	}

}

// To make sure MockRepository implement repository.Repository
var _ repository.Repository = &MockRepository{}
