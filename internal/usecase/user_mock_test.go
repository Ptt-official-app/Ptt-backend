package usecase

import (
	"context"
	"time"

	"github.com/PichuChen/go-bbs"
)

func (repo *MockRepository) GetUsers(ctx context.Context) ([]bbs.UserRecord, error) {
	ret := []bbs.UserRecord{}
	ret = append(ret, &MockUser{
		userId: "pichu",
	})
	return ret, nil
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
func (u *MockUser) NumBadPosts() int                     { return 0 }
func (u *MockUser) LastCountry() string                  { return "" }
func (u *MockUser) MailboxDescription() string           { return "" }
func (u *MockUser) ChessStatus() map[string]interface{}  { return map[string]interface{}{} }
func (u *MockUser) Plan() map[string]interface{}         { return map[string]interface{}{} }
