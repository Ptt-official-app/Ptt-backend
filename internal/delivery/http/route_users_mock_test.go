package http

import (
	"context"
	"fmt"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/go-bbs"
)

func (usecase *MockUsecase) GetUserByID(ctx context.Context, userID string) (bbs.UserRecord, error) {
	result := NewMockUserRecord(userID)
	return result, nil
}

func (usecase *MockUsecase) GetUserFavorites(ctx context.Context, userID string) ([]interface{}, error) {
	result := []interface{}{
		map[string]interface{}{
			"type":     "board",
			"board_id": "test_board_001",
		},
	}
	return result, nil
}

func (usecase *MockUsecase) GetUserInformation(ctx context.Context, userID string) (map[string]interface{}, error) {
	user := NewMockUserRecord("id")
	result := map[string]interface{}{
		"user_id":              user.UserID(),
		"nickname":             user.Nickname(),
		"realname":             user.RealName(),
		"number_of_login_days": fmt.Sprintf("%d", user.NumLoginDays()),
		"number_of_posts":      fmt.Sprintf("%d", user.NumPosts()),
		"number_of_badposts":   fmt.Sprintf("%d", user.NumBadPosts()),
		"money":                fmt.Sprintf("%d", user.Money()),
		"money_description":    "債台高築",
		"last_login_time":      user.LastLogin().Format(time.RFC3339),
		"last_login_ipv4":      user.LastHost(),
		"last_login_ip":        user.LastHost(),
		"last_login_country":   user.LastCountry(),
		"mailbox_description":  user.MailboxDescription(),
		"chess_status":         user.ChessStatus(),
		"plan":                 user.Plan(),
	}
	return result, nil
}

func (usecase *MockUsecase) GetUserPreferences(ctx context.Context, userID string) (map[string]string, error) {
	result := map[string]string{
		"favorite_no_highlight": "false",
	}

	return result, nil
}

func (usecase *MockUsecase) GetUserArticles(ctx context.Context, userID string) ([]interface{}, error) {
	return nil, nil
}

func (usecase *MockUsecase) GetUserComments(ctx context.Context, userID string) ([]interface{}, error) {
	result := []interface{}{
		map[string]interface{}{
			"board_id": "SYSOP",
		},
	}
	return result, nil
}

func (usecase *MockUsecase) GetUserDrafts(ctx context.Context, userID string, draftID string) (repository.UserDraft, error) {
	return nil, nil
}

func (usecase *MockUsecase) UpdateUserDraft(ctx context.Context, userID, draftID string, text []byte) (repository.UserDraft, error) {
	return nil, nil
}

func (usecase *MockUsecase) DeleteUserDraft(ctx context.Context, userID, draftID string) error {
	return nil
}

type MockUserRecord struct {
	userID             string
	nickname           string
	loginDays          int
	posts              int
	badPosts           int
	moneyDescription   string
	lastLoginTime      time.Time
	lastHost           string
	lastLoginCountry   string
	mailboxDescription string
	chessStatus        map[string]interface{}
	plan               map[string]interface{}
	realname           string
	money              int
}

// NewMockUserRecord generates fake user record for developing
func NewMockUserRecord(userID string) *MockUserRecord {
	return &MockUserRecord{
		userID:             userID,
		nickname:           "",
		loginDays:          0,
		posts:              0,
		badPosts:           0,
		moneyDescription:   "債台高築",
		lastLoginTime:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		lastHost:           "127.0.0.1",
		lastLoginCountry:   "",
		mailboxDescription: "",
		chessStatus:        map[string]interface{}{},
		plan:               map[string]interface{}{},
		realname:           "",
		money:              0,
	}
}
func (u *MockUserRecord) UserID() string { return u.userID }

// HashedPassword return user hashed password, it only for debug,
// If you want to check is user password correct, please use
// VerifyPassword insteaded.
func (u *MockUserRecord) HashedPassword() string { return "" }

// VerifyPassword will check user's password is OK. it will return null
// when OK and error when there are something wrong
func (u *MockUserRecord) VerifyPassword(password string) error { return nil }

// Nickname return a string for user's nickname, this string may change
// depend on user's mood, return empty string if this bbs system do not support
func (u *MockUserRecord) Nickname() string { return u.nickname }

// RealName return a string for user's real name, this string may not be changed
// return empty string if this bbs system do not support
func (u *MockUserRecord) RealName() string { return u.realname }

// NumLoginDays return how many days this have been login since account created.
func (u *MockUserRecord) NumLoginDays() int { return u.loginDays }

// NumPosts return how many posts this user has posted.
func (u *MockUserRecord) NumPosts() int { return u.posts }

// Money return the money this user have.
func (u *MockUserRecord) Money() int { return u.money }

// LastLogin return last login time of user
func (u *MockUserRecord) LastLogin() time.Time { return u.lastLoginTime }

// LastHost return last login host of user, it is IPv4 address usually, but it
// could be domain name or IPv6 address.
func (u *MockUserRecord) LastHost() string { return u.lastHost }

func (u *MockUserRecord) NumBadPosts() int { return u.badPosts }

func (u *MockUserRecord) LastCountry() string { return u.lastLoginCountry }

func (u *MockUserRecord) MailboxDescription() string { return u.mailboxDescription }

func (u *MockUserRecord) ChessStatus() map[string]interface{} { return u.chessStatus }

func (u *MockUserRecord) Plan() map[string]interface{} { return u.plan }
