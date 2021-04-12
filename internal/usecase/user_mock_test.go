package usecase

import (
	"context"
	"time"

	"github.com/Ptt-official-app/go-bbs"
)

func (repo *MockRepository) GetUsers(ctx context.Context) ([]bbs.UserRecord, error) {
	ret := []bbs.UserRecord{}
	ret = append(ret, &MockUser{
		userID: "pichu",
	})
	return ret, nil
}
func (repo *MockRepository) GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error) {
	return []bbs.FavoriteRecord{}, nil
}

func (repo *MockRepository) GetUserPreferences(ctx context.Context, userID string) (map[string]string, error) {
	return map[string]string{"favorite_no_highlight": "true"}, nil
}

func (repo *MockRepository) GetUserComments(ctx context.Context, userID string) ([]interface{}, error) {
	return []interface{}{
		map[string]interface{}{
			"board_id": "SYSOP",
		},
	}, nil
}

type MockUser struct {
	userID string
}

func (u *MockUser) UserID() string                       { return u.userID }
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
