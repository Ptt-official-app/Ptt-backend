package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

type userUsecase struct {
	logger   logging.Logger
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		logger:   logging.NewLogger(),
		userRepo: userRepo,
	}
}

func (u *userUsecase) GetUserByID(ctx context.Context, userID string) (bbs.UserRecord, error) {
	for _, it := range u.userRepo.GetUsers(ctx) {
		if userID == it.UserId() {
			return it, nil
		}
	}
	return nil, fmt.Errorf("user record not found")
}

func (u *userUsecase) GetUserFavorites(ctx context.Context, userID string) ([]interface{}, error) {
	recs, err := u.userRepo.GetUserFavoriteRecords(ctx, userID)
	if err != nil {
		return nil, err
	}

	dataItems := u.parseFavoriteFolderItem(recs)
	return dataItems, nil
}

func (u *userUsecase) GetUserInformation(ctx context.Context, userID string) (map[string]interface{}, error) {
	user, err := u.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get userrec for %s failed", userID)
	}

	// TODO: Check Etag or Not-Modified for cache

	result := map[string]interface{}{
		"user_id":              user.UserId(),
		"nickname":             user.Nickname(),
		"realname":             user.RealName(),
		"number_of_login_days": fmt.Sprintf("%d", user.NumLoginDays()),
		"number_of_posts":      fmt.Sprintf("%d", user.NumPosts()),
		// "number_of_badposts":   fmt.Sprintf("%d", user.NumLoginDays),
		"money":           fmt.Sprintf("%d", user.Money()),
		"last_login_time": user.LastLogin().Format(time.RFC3339),
		"last_login_ipv4": user.LastHost(),
		"last_login_ip":   user.LastHost(),
		// "last_login_country": fmt.Sprintf("%d", user.NumLoginDays),
		"chess_status": map[string]interface{}{},
		"plan":         map[string]interface{}{},
	}
	return result, nil
}

func (u *userUsecase) parseFavoriteFolderItem(recs []bbs.FavoriteRecord) []interface{} {
	dataItems := []interface{}{}
	for _, item := range recs {
		u.logger.Debugf("fav type: %v", item.Type())

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
				"items": u.parseFavoriteFolderItem(item.Records()),
			})

		case bbs.FavoriteTypeLine:
			dataItems = append(dataItems, map[string]interface{}{
				"type": "line",
			})
		default:
			u.logger.Warningf("parseFavoriteFolderItem unknown favItem type")
		}
	}
	return dataItems
}
