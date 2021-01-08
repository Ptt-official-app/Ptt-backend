package main

import (
	"encoding/json"
	// "fmt"
	"github.com/PichuChen/go-bbs"
	// "github.com/PichuChen/go-bbs/crypt"
	// "github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

func routeBoards(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("routeBoards: %v", r)

	// TODO: Check IP Flowspeed
	if r.Method == "GET" {
		getBoards(w, r)
		return
	}

}

func getBoards(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("getBoards: %v", r)

	token := getTokenFromRequest(r)
	userId, err := getUserIdFromToken(token)
	if err != nil {
		// user permission error
		// Support Guest?
		if !supportGuest() {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"token_invalid"}`))
			return
		} else {
			userId = "guest" // TODO: use const variable
		}
	}

	dataList := []interface{}{}
	for _, b := range boardHeader {
		// TODO: Show Board by user level
		if b.IsGroudBoard() {
			continue
		}
		if !shouldShowOnUserLevel(b, userId) {
			continue
		}
		jb, _ := json.Marshal(b)
		logger.Debugf("marshal board: %v", string(jb))
		dataList = append(dataList, map[string]interface{}{
			"id":             b.BrdName,
			"type":           "board",
			"title":          b.Title,
			"number_of_user": "0",
			"moderators":     strings.Split(b.BM, "/"), // TODO, set BM Split Token
		})
	}

	responseMap := map[string]interface{}{
		"data": dataList,
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)

}

func shouldShowOnUserLevel(b *bbs.BoardHeader, u string) bool {
	// TODO: Get user Level
	return true

}
