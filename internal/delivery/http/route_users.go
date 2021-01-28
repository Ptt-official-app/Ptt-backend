package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func (delivery *httpDelivery) routeUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed
	switch r.Method {
	case http.MethodGet:
		delivery.getUsers(w, r)
	}
}

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
	err := checkTokenPermission(token,
		[]permission{PermissionReadUserInformation},
		map[string]string{
			"user_id": userId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataMap, err := delivery.userUsecase.GetUserInformation(context.Background(), userId)
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
	err := checkTokenPermission(token,
		[]permission{PermissionReadUserInformation},
		map[string]string{
			"user_id": userId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataItems, err := delivery.userUsecase.GetUserFavorites(context.Background(), userId)
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

func parseUserPath(path string) (userId string, item string, err error) {
	pathSegment := strings.Split(path, "/")
	// /{{version}}/users/{{user_id}}/{{item}}
	if len(pathSegment) == 4 {
		// /{{version}}/users/{{user_id}}
		return pathSegment[3], "", nil
	}
	return pathSegment[3], pathSegment[4], nil
}
