package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-bbs"
)

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
func (u *MockUser) UserFlag() uint32                     { return 0x02000A60 }

func TestGetUsers(t *testing.T) {

	repo := repository{
		userRecords:  []bbs.UserRecord{},
		boardRecords: []bbs.BoardRecord{},
	}

	actual, err := repo.GetUsers(context.TODO())
	if err != nil {
		t.Errorf("GetUsers excepted err == nil, got %v", err)
		return
	}

	if actual == nil {
		t.Errorf("GetUsers got %v, expected not equal nil", actual)
	}

}

func TestGetUserArticles(t *testing.T) {

	repo := repository{
		userRecords: []bbs.UserRecord{
			&MockUser{
				userID: "SYSOP",
			},
		},
		boardRecords: []bbs.BoardRecord{},
	}

	actual, err := repo.GetUserPreferences(context.TODO(), "SYSOP")
	if err != nil {
		t.Errorf("GetUserPreferences excepted err == nil, got %v", err)
		return
	}
	t.Logf("%v", actual)
	if actual == nil {
		t.Errorf("GetUserPreferences got %v, expected not equal nil", actual)
	}
	if actual["board_sort"] != "true" {
		t.Errorf("GetUserPreferences boart_sort got %v, expected true", actual["board_sort"])
	}

}
