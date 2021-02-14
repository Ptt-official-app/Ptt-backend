package usecase

import (
	"context"
	"github.com/PichuChen/go-bbs"
	"time"
)

func (repo *MockRepository) GetUsers(ctx context.Context) []bbs.UserRecord {
	ret := []bbs.UserRecord{}
	ret = append(ret, &MockUser{
		userId: "pichu",
	})
	return ret
}
func (repo *MockRepository) GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error) {
	return []bbs.FavoriteRecord{}, nil
}

type MockUser struct {
	userId string
}

func (u *MockUser) UserId() string                       { return u.userId }
func (u *MockUser) HashedPassword() string               { return "" }
func (u *MockUser) VerifyPassword(password string) error { return nil }
func (u *MockUser) Nickname() string                     { return "" }
func (u *MockUser) RealName() string                     { return "" }
func (u *MockUser) NumLoginDays() int                    { return 0 }
func (u *MockUser) NumPosts() int                        { return 0 }
func (u *MockUser) Money() int                           { return 0 }
func (u *MockUser) LastLogin() time.Time                 { return time.Unix(0, 0) }
func (u *MockUser) LastHost() string                     { return "" }
