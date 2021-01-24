package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func (rest *restHandler) routeUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed
	switch r.Method {
	case http.MethodGet:
		rest.getUsers(w, r)
	}
}

func (rest *restHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	userId, item, err := parseUserPath(r.URL.Path)
	switch item {
	case "information":
		rest.getUserInformation(w, r, userId)
	case "favorites":
		rest.getUserFavorites(w, r, userId)
	default:
		rest.logger.Noticef("user id: %v not exist but be queried, info: %v err: %v", userId, item, err)
		w.WriteHeader(http.StatusNotFound)
	}
}

func (rest *restHandler) getUserInformation(w http.ResponseWriter, r *http.Request, userId string) {
	token := rest.getTokenFromRequest(r)
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

	dataMap, err := rest.userUsecase.GetUserInformation(context.Background(), userId)
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

func (rest *restHandler) getUserFavorites(w http.ResponseWriter, r *http.Request, userId string) {
	token := rest.getTokenFromRequest(r)
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

	dataItems, err := rest.userUsecase.GetUserFavorites(context.Background(), userId)
	if err != nil {
		rest.logger.Errorf("failed to get user favorites: %s\n", err)
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
