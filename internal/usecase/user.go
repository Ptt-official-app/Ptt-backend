package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/go-bbs"
)

func (usecase *usecase) GetUserByID(ctx context.Context, userID string) (bbs.UserRecord, error) {
	users, err := usecase.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, it := range users {
		if userID == it.UserID() {
			return it, nil
		}
	}
	return nil, fmt.Errorf("user record not found")
}

func (usecase *usecase) GetUserFavorites(ctx context.Context, userID string) ([]interface{}, error) {
	recs, err := usecase.repo.GetUserFavoriteRecords(ctx, userID)
	if err != nil {
		return nil, err
	}

	dataItems := usecase.parseFavoriteFolderItem(recs)
	return dataItems, nil
}

func (usecase *usecase) GetUserInformation(ctx context.Context, userID string) (map[string]interface{}, error) {
	userrec, err := usecase.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get userrec for %s failed", userID)
	}
	user := userrec.(repository.BBSUserRecord)

	// TODO: Check Etag or Not-Modified for cache

	result := map[string]interface{}{
		"user_id":              user.UserID(),
		"nickname":             user.Nickname(),
		"realname":             user.RealName(),
		"number_of_login_days": strconv.FormatInt(int64(user.NumLoginDays()), 10),
		"number_of_posts":      strconv.FormatInt(int64(user.NumPosts()), 10),
		"number_of_badposts":   strconv.FormatInt(int64(user.NumBadPosts()), 10),
		"money":                strconv.FormatInt(int64(user.Money()), 10),
		"money_description":    getMoneyDiscription(user.Money()),
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

func (usecase *usecase) GetUserArticles(ctx context.Context, userID string) ([]interface{}, error) {
	dataItems := []interface{}{}

	// Because there is no user’s historical article data stored, first get all public boards, and then get user articles
	boards := usecase.GetBoards(ctx, userID)

	for _, board := range boards {
		articleRecords, err := usecase.repo.GetUserArticles(ctx, board.BoardID())
		if err != nil {
			return nil, err
		}

		for index := range articleRecords {
			if articleRecords[index].Owner() == userID {
				dataItems = append(dataItems, map[string]interface{}{
					"board_id":        "", // FIXME: use concrete value rather than ""
					"filename":        articleRecords[index].Filename(),
					"modified_time":   articleRecords[index].Modified(),
					"recommend_count": articleRecords[index].Recommend(),
					"comment_count":   0,  // FIXME: use concrete value rather than 0
					"post_date":       "", // FIXME: use concrete value rather than ""
					"title":           articleRecords[index].Title(),
					"money":           articleRecords[index].Money(),
					"owner":           articleRecords[index].Owner(),
					"aid":             "", // FIXME: use concrete value rather than ""
					"url":             "", // FIXME: use concrete value rather than ""
				})
			}
		}
	}

	return dataItems, nil
}

func (usecase *usecase) GetUserPreferences(ctx context.Context, userID string) (map[string]string, error) {
	rec, err := usecase.repo.GetUserPreferences(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get preferences_rec for %s failed", userID)
	}

	return rec, nil
}

func (usecase *usecase) GetUserComments(ctx context.Context, userID string) ([]interface{}, error) {
	dataItems, err := usecase.repo.GetUserComments(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dataItems, nil
}

func (usecase *usecase) GetUserDrafts(ctx context.Context, userID, draftID string) ([]byte, error) {
	if !isValidDraftID([]byte(draftID)) {
		return nil, fmt.Errorf("invalid draft ID: %s", draftID)
	}
	return usecase.repo.GetUserDrafts(ctx, userID, draftID)
}

func (usecase *usecase) UpdateUserDraft(ctx context.Context, userID, draftID string, text []byte) ([]byte, error) {
	if !isValidDraftID([]byte(draftID)) {
		return nil, fmt.Errorf("invalid draft ID: %s", draftID)
	}
	return usecase.repo.UpdateUserDraft(ctx, userID, draftID, text)
}

func (usecase *usecase) DeleteUserDraft(ctx context.Context, userID, draftID string) error {
	if !isValidDraftID([]byte(draftID)) {
		return fmt.Errorf("invalid draft ID: %s", draftID)
	}
	return usecase.repo.DeleteUserDraft(ctx, userID, draftID)
}

func (usecase *usecase) parseFavoriteFolderItem(recs []bbs.FavoriteRecord) []interface{} {
	dataItems := []interface{}{}
	for _, item := range recs {
		usecase.logger.Debugf("fav type: %v", item.Type())

		switch item.Type() {
		case bbs.FavoriteTypeBoard:
			dataItems = append(dataItems, map[string]interface{}{
				"type":     "board",
				"board_id": item.BoardID(),
			})

		case bbs.FavoriteTypeFolder:
			dataItems = append(dataItems, map[string]interface{}{
				"type":  "folder",
				"title": item.Title(),
				"items": usecase.parseFavoriteFolderItem(item.Records()),
			})

		case bbs.FavoriteTypeLine:
			dataItems = append(dataItems, map[string]interface{}{
				"type": "line",
			})
		default:
			usecase.logger.Warningf("parseFavoriteFolderItem unknown favItem type")
		}
	}
	return dataItems
}

func isValidDraftID(draftID []byte) bool {
	if len(draftID) == 1 {
		return draftID[0] >= '0' && draftID[0] <= '9'
	}
	return false
}

func getMoneyDiscription(money int) string {
	var result string
	switch {
	case money <= 9:
		result = "債台高築"
	case money <= 109:
		result = "赤貧"
	case money <= 1099:
		result = "清寒"
	case money <= 10999:
		result = "普通"
	case money <= 109999:
		result = "小康"
	case money <= 1099999:
		result = "小富"
	case money <= 10999999:
		result = "中富"
	case money <= 109999999:
		result = "大富翁"
	case money <= 1099999999:
		result = "富可敵國"
	default:
		result = "比爾蓋天"
	}
	return result
}
