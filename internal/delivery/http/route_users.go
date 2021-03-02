package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

// ref - internal\delivery\http\route.go:52
// route map
func (delivery *Delivery) getUsers(w http.ResponseWriter, r *http.Request) {
	userID, item, err := parseUserPath(r.URL.Path)
	switch item {
	case "information":
		delivery.getUserInformation(w, r, userID)
	case "favorites":
		delivery.getUserFavorites(w, r, userID)
	default:
		delivery.logger.Noticef("user id: %v not exist but be queried, info: %v err: %v", userID, item, err)
		w.WriteHeader(http.StatusNotFound)
	}
}

// ref - https://pttapp.cc/swagger/#/%E4%BD%BF%E7%94%A8%E8%80%85%E9%83%A8%E5%88%86/get_v1_users__user_id__information
// desc - 取得使用者資訊，包括上次上站位置等
// route - /v1/users/{user_id}/information
// parameter - userID : 使用者 id
func (delivery *Delivery) getUserInformation(w http.ResponseWriter, r *http.Request, userID string) {
	token := delivery.getTokenFromRequest(r)

	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadUserInformation},
		map[string]string{
			"user_id": userID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataMap, err := delivery.usecase.GetUserInformation(context.Background(), userID)
	if err != nil {
		// TODO: record error
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_userrec_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getUserInformation error response err: %w", err)
		}
		return
	}

	responseMap := map[string]interface{}{
		"data": dataMap,
	}
	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	_, err = w.Write(responseByte)
	if err != nil {
		delivery.logger.Errorf("getUserInformation success response err: %w", err)
	}
}

// ref - https://pttapp.cc/swagger/#/%E4%BD%BF%E7%94%A8%E8%80%85%E9%83%A8%E5%88%86/get_v1_users__user_id__favorites
// desc - 我的最愛列表
// route - /v1/users/{user_id}/favorites
// parameter - userID : 使用者 id
func (delivery *Delivery) getUserFavorites(w http.ResponseWriter, r *http.Request, userID string) {
	token := delivery.getTokenFromRequest(r)
	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadUserInformation},
		map[string]string{
			"user_id": userID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataItems, err := delivery.usecase.GetUserFavorites(context.Background(), userID)
	if err != nil {
		delivery.logger.Errorf("failed to get user favorites: %s\n", err)
	}

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": dataItems,
		},
	}

	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	_, err = w.Write(responseByte)
	if err != nil {
		delivery.logger.Errorf("getUserFavorites success response err: %w", err)
	}
}
