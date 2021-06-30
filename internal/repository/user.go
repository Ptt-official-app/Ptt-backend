package repository

import (
	"context"
	"fmt"

	"github.com/Ptt-official-app/go-bbs"
	// TODO: remove direct access pttbbs, implement it in go-bbs package
	"github.com/Ptt-official-app/go-bbs/pttbbs"
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

	var u bbs.UserRecord = nil
	for _, it := range repo.userRecords {
		if it.UserID() == userID {
			u = it
		}
	}

	if u == nil {
		// not found
		return nil, fmt.Errorf("userrecord not found")
	}

	c := u.UserFlag()
	logger.Debugf("userflag of %s is %X", userID, c)

	rawRec, isPttbbs := u.(*pttbbs.Userec)
	pagerUIType := "origin"
	if isPttbbs {
		switch rawRec.PagerUIType {
		case 0:
			pagerUIType = "origin"
		case 1:
			pagerUIType = "new"
		case 2:
			pagerUIType = "ofo"
		default:
			break
		}
	}

	result := map[string]string{
		"favorite_no_highlight":      fmt.Sprintf("%v", c&pttbbs.UfFavNohilight != 0),
		"favorite_add_new":           fmt.Sprintf("%v", c&pttbbs.UfFavAddnew != 0),
		"friend":                     fmt.Sprintf("%v", c&pttbbs.UfFriend != 0),
		"board_sort":                 fmt.Sprintf("%v", c&pttbbs.UfBrdsort != 0),
		"ad_banner":                  fmt.Sprintf("%v", c&pttbbs.UfAdbanner != 0),
		"ad_banner_user_song":        fmt.Sprintf("%v", c&pttbbs.UfAdbannerUsong != 0),
		"dbcs_aware":                 fmt.Sprintf("%v", c&pttbbs.UfDbcsAware != 0),
		"dbcs_no_interupting_escape": fmt.Sprintf("%v", c&pttbbs.UfDbcsNointresc != 0),
		"dbcs_drop_repeat":           fmt.Sprintf("%v", c&pttbbs.UfDbscDropRepeat != 0),
		"no_modification_mark":       fmt.Sprintf("%v", c&pttbbs.UfNoModmark != 0),
		"colored_modification_mark":  fmt.Sprintf("%v", c&pttbbs.UfColoredModmark != 0),
		"default_backup":             fmt.Sprintf("%v", c&pttbbs.UfDefbackup != 0),
		"new_angel_pager":            fmt.Sprintf("%v", c&pttbbs.UfNewAngelPager != 0),
		"reject_outside_mail":        fmt.Sprintf("%v", c&pttbbs.UfRejOuttamail != 0),
		"secure_login":               fmt.Sprintf("%v", c&pttbbs.UfSecureLogin != 0),
		"foreign":                    fmt.Sprintf("%v", c&pttbbs.UfForeign != 0),
		"live_right":                 fmt.Sprintf("%v", c&pttbbs.UfLiveright != 0),
		"menu_lightbar":              fmt.Sprintf("%v", c&pttbbs.UfMenuLightbar != 0),
		"cursor_ascii":               fmt.Sprintf("%v", c&pttbbs.UfCursorASCII != 0),
		"pager_ui":                   pagerUIType,
	}

	return result, nil
}

// TODO: no required method in go-bbs and we use a mock, replace it when available
func (repo *repository) GetUserComments(_ context.Context, userID string) ([]bbs.UserCommentRecord, error) {
	result, err := repo.db.GetUserCommentRecordFile(userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *repository) GetUserDrafts(_ context.Context, userID, draftID string) (bbs.UserDraft, error) {
	return repo.db.GetUserDrafts(userID, draftID)
}

func (repo *repository) UpdateUserDraft(_ context.Context, userID, draftID string, text []byte) (bbs.UserDraft, error) {
	return nil, nil
}

func (repo *repository) DeleteUserDraft(_ context.Context, userID, draftID string) error {
	return nil
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
