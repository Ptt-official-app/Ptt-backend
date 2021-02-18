package repository

import (
	"context"
	"fmt"

	"github.com/PichuChen/go-bbs"
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

func loadUserRecords(db *bbs.DB) ([]bbs.UserRecord, error) {
	userRecords, err := db.ReadUserRecords()
	if err != nil {
		logger.Errorf("get user rec error: %v", err)
		return nil, fmt.Errorf("failed to read user records: %w", err)
	}
	results := make([] bbs.UserRecord, 0, len(userRecords))
	for _, rec := range userRecords {
		results = append(results, &bbsUserRecord{rec})
	}
	return results, nil
}
