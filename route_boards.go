package main

import (
	"encoding/json"
	// "fmt"
	"github.com/PichuChen/go-bbs"
	// "github.com/PichuChen/go-bbs/crypt"
	// "github.com/dgrijalva/jwt-go"
	"fmt"
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
	boardId, item, err := parseBoardPath(r.URL.Path)
	if boardId == "" {
		getBoardList(w, r)
		return
	}
	// get single board
	if item == "information" {
		getBoardInformation(w, r, boardId)
		return
	}

	logger.Noticef("board id: %v not exist but be queried, info: %v err: %v", boardId, item, err)
}

func getBoardList(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("getBoardList: %v", r)

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
		dataList = append(dataList, marshalBoardHeader(b))
	}

	responseMap := map[string]interface{}{
		"data": dataList,
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)

}

func getBoardInformation(w http.ResponseWriter, r *http.Request, boardId string) {
	logger.Debugf("getBoardInformation: %v", r)
	token := getTokenFromRequest(r)
	err := checkTokenPermission(token,
		[]permission{PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	brd, err := findBoardHeaderById(boardId)
	if err != nil {
		// TODO: record error
		logger.Warningf("find board %s failed: %v", boardId, err)
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_board_error",
			"error_description": "get board for " + boardId + " failed",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		w.Write(b)
		return
	}

	responseMap := map[string]interface{}{
		"data": marshalBoardHeader(brd),
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)

}

func marshalBoardHeader(b *bbs.BoardHeader) map[string]interface{} {
	return map[string]interface{}{
		"id":             b.BrdName,
		"type":           "board",
		"title":          b.Title,
		"number_of_user": "0",
		"moderators":     strings.Split(b.BM, "/"), // TODO, set BM Split Token
	}

}

func shouldShowOnUserLevel(b *bbs.BoardHeader, u string) bool {
	// TODO: Get user Level
	return true

}

// parseBoardPath covert url path from /v1/boards/SYSOP/article to
// {SYSOP, article) or /v1/boards to {,}
func parseBoardPath(path string) (boardId string, item string, err error) {
	pathSegment := strings.Split(path, "/")
	if len(pathSegment) == 5 {
		// /{{version}}/boards/{{board_id}}/{{item}}
		return pathSegment[3], pathSegment[4], nil
	} else if len(pathSegment) == 4 {
		// /{{version}}/boards/{{board_id}}
		return pathSegment[3], "", nil
	} else if len(pathSegment) == 3 {
		// /{{version}}/boards
		return "", "", nil
	}
	logger.Warningf("parseBoardPath got malform path: %v", path)
	return "", "", nil

}

func findBoardHeaderById(boardId string) (*bbs.BoardHeader, error) {

	for _, it := range boardHeader {
		if boardId == it.BrdName {
			return it, nil
		}
	}
	return nil, fmt.Errorf("board record not found")

}
