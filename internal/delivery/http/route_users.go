package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (delivery *httpDelivery) getUsers(w http.ResponseWriter, r *http.Request) {
	userId, item, err := parseUserPath(r.URL.Path)
	switch item {
	case "information":
		delivery.getUserInformation(w, r, userId)
	case "favorites":
		delivery.getUserFavorites(w, r, userId)
	default:
		delivery.logger.Noticef("user id: %v not exist but be queried, info: %v err: %v", userId, item, err)
		w.WriteHeader(http.StatusNotFound)
	}
}

func (delivery *httpDelivery) getUserInformation(w http.ResponseWriter, r *http.Request, userId string) {
	token := delivery.getTokenFromRequest(r)

	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadUserInformation},
		map[string]string{
			"user_id": userId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataMap, err := delivery.usecase.GetUserInformation(context.Background(), userId)
	if err != nil {
		// TODO: record error
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_userrec_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		w.Write(b)
		return
	}

	responseMap := map[string]interface{}{
		"data": dataMap,
	}
	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	w.Write(responseByte)
}

func (delivery *httpDelivery) getUserFavorites(w http.ResponseWriter, r *http.Request, userId string) {
	token := delivery.getTokenFromRequest(r)
	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadUserInformation},
		map[string]string{
			"user_id": userId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataItems, err := delivery.usecase.GetUserFavorites(context.Background(), userId)
	if err != nil {
		delivery.logger.Errorf("failed to get user favorites: %s\n", err)
	}

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": dataItems,
		},
	}

	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	w.Write(responseByte)
}
