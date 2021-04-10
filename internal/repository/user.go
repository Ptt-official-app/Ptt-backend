package repository

import (
	"context"
	"fmt"

	"github.com/Ptt-official-app/go-bbs"
)

// BBSUserRecord : currently interface `bbs.UserRecord` of go-bbs
// lacks some required methods, so we use a mock
type BBSUserRecord interface {
	bbs.UserRecord
	NumBadPosts() int
	LastCountry() string
	MailboxDescription() string
	ChessStatus() map[string]interface{}
	Plan() map[string]interface{}
}

type bbsUserRecord struct {
	bbs.UserRecord
}

// NumBadPosts returns how many bad posts this user has posted.
func (record *bbsUserRecord) NumBadPosts() int {
	return 0
}

// LastCountry returns last login country of user
func (record *bbsUserRecord) LastCountry() string {
	return ""
}

// MailboxDescription returns the mailbox description of user
func (record *bbsUserRecord) MailboxDescription() string {
	return ""
}

// ChessStatus returns chess status
func (record *bbsUserRecord) ChessStatus() map[string]interface{} {
	return map[string]interface{}{}
}

// Plan returns plan
func (record *bbsUserRecord) Plan() map[string]interface{} {
	return map[string]interface{}{}
}

func (repo *repository) GetUsers(_ context.Context) ([]bbs.UserRecord, error) {
	return repo.userRecords, nil
}

func (repo *repository) GetUserFavoriteRecords(ctx context.Context, userID string) ([]bbs.FavoriteRecord, error) {
	return repo.db.ReadUserFavoriteRecords(userID)
}

func (repo *repository) GetUserArticles(_ context.Context, boardID string) ([]bbs.ArticleRecord, error) {
	return repo.db.ReadBoardArticleRecordsFile(boardID)
}

// TODO: no required method in go-bbs and we use a mock, replace it when available
func (repo *repository) GetUserPreferences(_ context.Context, userID string) (map[string]string, error) {
	result := map[string]string{
		"favorite_no_highlight":      "No value",
		"favorite_add_new":           "No value",
		"friend":                     "No value",
		"board_sort":                 "No value",
		"ad_banner":                  "No value",
		"ad_banner_user_song":        "No value",
		"dbcs_aware":                 "No value",
		"dbcs_no_interupting_escape": "No value",
		"dbcs_drop_repeat":           "No value",
		"no_modification_mark":       "No value",
		"colored_modification_mark":  "No value",
		"default_backup":             "No value",
		"new_angel_pager":            "No value",
		"reject_outside_mail":        "No value",
		"secure_login":               "No value",
		"foreign":                    "No value",
		"live_right":                 "No value",
		"menu_lightbar":              "No value",
		"cursor_ascii":               "No value",
		"pager_ui":                   "No value",
	}
	return result, nil
}

// TODO: no required method in go-bbs and we use a mock, replace it when available
func (repo *repository) GetUserComments(_ context.Context, userID string) ([]interface{}, error) {
	result := []interface{}{
		map[string]interface{}{
			"board_id":        "No value",
			"filename":        "No value",
			"modified_time":   "No value",
			"recommend_count": "No value",
			"comment_count":   "No value",
			"post_data":       "No value",
			"title":           "No value",
			"url":             "No value",
			"comment_order":   "No value",
			"comment_time":    "No value",
		},
	}
	return result, nil
}

func loadUserRecords(db *bbs.DB) ([]bbs.UserRecord, error) {
	userRecords, err := db.ReadUserRecords()
	if err != nil {
		logger.Errorf("get user rec error: %v", err)
		return nil, fmt.Errorf("failed to read user records: %w", err)
	}
	results := make([]bbs.UserRecord, 0, len(userRecords))
	for _, rec := range userRecords {
		results = append(results, &bbsUserRecord{rec})
	}
	return results, nil
}
