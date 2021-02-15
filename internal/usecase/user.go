package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

func (usecase *usecase) GetUserByID(ctx context.Context, userID string) (bbs.UserRecord, error) {
	users, err := usecase.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, it := range users {
		if userID == it.UserId() {
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
		"user_id":              user.UserId(),
		"nickname":             user.Nickname(),
		"realname":             user.RealName(),
		"number_of_login_days": fmt.Sprintf("%d", user.NumLoginDays()),
		"number_of_posts":      fmt.Sprintf("%d", user.NumPosts()),
		"number_of_badposts":   fmt.Sprintf("%d", user.NumBadPosts()),
		"money":           fmt.Sprintf("%d", user.Money()),
		"money_description": getMoneyDiscription(user.Money()),
		"last_login_time": user.LastLogin().Format(time.RFC3339),
		"last_login_ipv4": user.LastHost(),
		"last_login_ip":   user.LastHost(),
		"last_login_country": user.LastCountry(),
		"mailbox_description": user.MailboxDescription(),
		"chess_status": user.ChessStatus(),
		"plan":         user.Plan(),
	}
	return result, nil
}

func (usecase *usecase) parseFavoriteFolderItem(recs []bbs.FavoriteRecord) []interface{} {
	dataItems := []interface{}{}
	for _, item := range recs {
		usecase.logger.Debugf("fav type: %v", item.Type())

		switch item.Type() {
		case bbs.FavoriteTypeBoard:
			dataItems = append(dataItems, map[string]interface{}{
				"type":     "board",
				"board_id": item.BoardId(),
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
