package main

import (
	"github.com/PichuChen/go-bbs"

	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// routeBoards is the handler for `/v1/boards`
func routeBoards(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("routeBoards: %v", r)

	// TODO: Check IP Flowspeed
	if r.Method == http.MethodGet {
		getBoards(w, r)
		return
	}
}

// getBoards is the handler for `/v1/boards` with GET method
func getBoards(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("getBoards: %v", r)
	boardID, item, filename, err := parseBoardPath(r.URL.Path)

	if boardID == "" {
		getBoardList(w, r)
		return
	}
	// get single board
	switch item {
	case "information":
		getBoardInformation(w, r, boardID)
	case "articles":
		if filename == "" {
			getBoardArticles(w, r, boardID)
		} else {
			getBoardArticlesFile(w, r, boardID, filename)
		}
	case "treasures":
		getBoardTreasures(w, r, boardID)
	default:
		// 404
		w.WriteHeader(http.StatusNotFound)
		logger.Noticef("board id: %v not exist but be queried, info: %v err: %v", boardID, item, err)
	}
}

func getBoardList(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("getBoardList: %v", r)

	token := getTokenFromRequest(r)
	userID, err := getUserIDFromToken(token)

	if err != nil {
		// user permission error
		// Support Guest?
		if !supportGuest() {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"token_invalid"}`))

			return
		}

		userID = "guest" // TODO: use const variable
	}

	dataList := []interface{}{}

	for _, b := range boardHeader {
		// TODO: Show Board by user level
		if b.IsClass() {
			continue
		}

		if !shouldShowOnUserLevel(b, userID) {
			continue
		}

		jb, err := json.Marshal(b)
		if err != nil {
			logger.Warningf("failed to marshal board: %s\n", err)
		}

		logger.Debugf("marshal board: %v", string(jb))

		dataList = append(dataList, marshalBoardHeader(b))
	}

	responseMap := map[string]interface{}{
		"data": dataList,
	}

	b, err := json.MarshalIndent(responseMap, "", "  ")
	if err != nil {
		logger.Errorf("failed to marshal response data: %s\n", err)
	}

	if _, err := w.Write(b); err != nil {
		logger.Errorf("failed to write response: %s\n", err)
	}
}

func getBoardInformation(w http.ResponseWriter, r *http.Request, boardID string) {
	logger.Debugf("getBoardInformation: %v", r)
	token := getTokenFromRequest(r)
	err := checkTokenPermission(token,
		[]permission{PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	brd, err := findBoardHeaderByID(boardID)
	if err != nil {
		// TODO: record error
		logger.Warningf("find board %s failed: %v", boardID, err)
		w.WriteHeader(http.StatusInternalServerError)

		m := map[string]string{
			"error":             "find_board_error",
			"error_description": "get board for " + boardID + " failed",
		}
		b, err := json.MarshalIndent(m, "", "  ")

		if err != nil {
			logger.Errorf("failed to marshal response data: %s\n", err)
		}

		if _, err := w.Write(b); err != nil {
			logger.Errorf("failed to write response: %s\n", err)
		}

		return
	}

	responseMap := map[string]interface{}{
		"data": marshalBoardHeader(brd),
	}

	b, err := json.MarshalIndent(responseMap, "", "  ")
	if err != nil {
		logger.Errorf("failed to marshal response data: %v\n", err)
	}

	if _, err := w.Write(b); err != nil {
		logger.Errorf("failed to write response: %v\n", err)
	}
}

// marshal generate board or class metadata object,
// b is input header
func marshalBoardHeader(b bbs.BoardRecord) map[string]interface{} {
	ret := map[string]interface{}{
		"title":          b.Title(),
		"number_of_user": "0",
		"moderators":     b.BM(),
	}
	if b.IsClass() {
		// class
		// Assign ID from foreach loop
		ret["type"] = "class"
	} else {
		// board
		ret["id"] = b.BoardId()
		ret["type"] = "board"
	}

	return ret
}

func shouldShowOnUserLevel(b bbs.BoardRecord, u string) bool {
	// TODO: Get user Level
	return true
}

// parseBoardPath covert url path from /v1/boards/SYSOP/article to
// {SYSOP, article) or /v1/boards to {,}
func parseBoardPath(path string) (boardID string, item string, filename string, err error) {
	pathSegment := strings.Split(path, "/")

	if len(pathSegment) >= 6 {
		// /{{version}}/boards/{{class_id}}/{{item}}/{{filename}}
		boardID = pathSegment[3]
		item = pathSegment[4]
		filename = pathSegment[5]

		return
	} else if len(pathSegment) == 5 {
		// /{{version}}/boards/{{class_id}}/{{item}}
		boardID = pathSegment[3]
		item = pathSegment[4]

		return
	} else if len(pathSegment) == 4 {
		// /{{version}}/boards/{{class_id}}
		boardID = pathSegment[3]

		return
	} else if len(pathSegment) == 3 {
		// /{{version}}/boards
		// Should not be reach...

		return
	}

	logger.Warningf("parseBoardPath got malformed path: %v", path)

	return
}

func findBoardHeaderByID(boardID string) (bbs.BoardRecord, error) {
	for _, it := range boardHeader {
		if boardID == it.BoardId() {
			return it, nil
		}
	}

	return nil, fmt.Errorf("board record not found")
}
