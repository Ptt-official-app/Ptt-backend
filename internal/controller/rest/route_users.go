package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PichuChen/go-bbs"
)

func (rest *restHandler) routeUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed

	if r.Method == "GET" {
		rest.getUsers(w, r)
		return
	}

}

func (rest *restHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	userId, item, err := parseUserPath(r.URL.Path)

	if item == "information" {
		rest.getUserInformation(w, r, userId)
		return
	} else if item == "favorites" {
		rest.getUserFavorites(w, r, userId)
		return
	}
	// else
	rest.logger.Noticef("user id: %v not exist but be queried, info: %v err: %v", userId, item, err)
	w.WriteHeader(http.StatusNotFound)
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

	userrec, err := rest.findUserecById(userId)
	if err != nil {
		// TODO: record error

		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_userrec_error",
			"error_description": "get userrec for " + userId + " failed",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		w.Write(b)
		return
	}

	// TODO: Check Etag or Not-Modified for cache

	dataMap := map[string]interface{}{
		"user_id":              userrec.UserId(),
		"nickname":             userrec.Nickname(),
		"realname":             userrec.RealName(),
		"number_of_login_days": fmt.Sprintf("%d", userrec.NumLoginDays()),
		"number_of_posts":      fmt.Sprintf("%d", userrec.NumPosts()),
		// "number_of_badposts":   fmt.Sprintf("%d", userrec.NumLoginDays),
		"money":           fmt.Sprintf("%d", userrec.Money()),
		"last_login_time": userrec.LastLogin().Format(time.RFC3339),
		"last_login_ipv4": userrec.LastHost(),
		"last_login_ip":   userrec.LastHost(),
		// "last_login_country": fmt.Sprintf("%d", userrec.NumLoginDays),
		"chess_status": map[string]interface{}{},
		"plan":         map[string]interface{}{},
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

	recs, err := rest.userRepo.GetUserFavoriteRecords(context.Background(), userId)
	rest.logger.Debugf("file items length: %v", len(recs))
	// dataMap := map[string]interface{}{}

	dataItems := rest.parseFavoriteFolderItem(recs)

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": dataItems,
		},
	}

	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	w.Write(responseByte)
}

func (rest *restHandler) parseFavoriteFolderItem(recs []bbs.FavoriteRecord) []interface{} {
	dataItems := []interface{}{}
	for _, item := range recs {
		rest.logger.Debugf("fav type: %v", item.Type())

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
				"items": rest.parseFavoriteFolderItem(item.Records()),
			})

		case bbs.FavoriteTypeLine:
			dataItems = append(dataItems, map[string]interface{}{
				"type": "line",
			})
		default:
			rest.logger.Warningf("parseFavoriteFolderItem unknown favItem type")
		}
	}
	return dataItems
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
