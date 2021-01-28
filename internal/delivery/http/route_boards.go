package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/PichuChen/go-bbs"
)

// routeBoards is the handler for `/v1/boards`
func (delivery *httpDelivery) routeBoards(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("routeBoards: %v", r)

	// TODO: Check IP Flowspeed
	if r.Method == "GET" {
		delivery.getBoards(w, r)
		return
	}

}

// getBoards is the handler for `/v1/boards` with GET method
func (delivery *httpDelivery) getBoards(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getBoards: %v", r)
	boardId, item, filename, err := delivery.parseBoardPath(r.URL.Path)
	if boardId == "" {
		delivery.getBoardList(w, r)
		return
	}
	// get single board
	if item == "information" {
		delivery.getBoardInformation(w, r, boardId)
		return
	} else if item == "articles" {
		if filename == "" {
			delivery.getBoardArticles(w, r, boardId)
		} else {
			delivery.getBoardArticlesFile(w, r, boardId, filename)
		}
		return
	} else if item == "treasures" {
		delivery.getBoardTreasures(w, r, boardId)
		return
	}

	// 404
	w.WriteHeader(http.StatusNotFound)

	delivery.logger.Noticef("board id: %v not exist but be queried, info: %v err: %v", boardId, item, err)
}

func (delivery *httpDelivery) getBoardList(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getBoardList: %v", r)

	token := delivery.getTokenFromRequest(r)
	userId, err := delivery.getUserIdFromToken(token)
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

	boards := delivery.boardUsecase.GetBoards(context.Background(), userId)

	dataList := make([]interface{}, 0, len(boards))
	for _, board := range boards {
		dataList = append(dataList, marshalBoardHeader(board))
	}

	responseMap := map[string]interface{}{
		"data": dataList,
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)
}

func (delivery *httpDelivery) getBoardInformation(w http.ResponseWriter, r *http.Request, boardId string) {
	delivery.logger.Debugf("getBoardInformation: %v", r)
	token := delivery.getTokenFromRequest(r)
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

	brd, err := delivery.boardUsecase.GetBoardByID(context.Background(), boardId)
	if err != nil {
		// TODO: record error
		delivery.logger.Warningf("find board %s failed: %v", boardId, err)
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
func (delivery *httpDelivery) parseBoardPath(path string) (boardId string, item string, filename string, err error) {
	pathSegment := strings.Split(path, "/")

	if len(pathSegment) >= 6 {
		// /{{version}}/boards/{{class_id}}/{{item}}/{{filename}}
		boardId = pathSegment[3]
		item = pathSegment[4]
		filename = pathSegment[5]
		return
	} else if len(pathSegment) == 5 {
		// /{{version}}/boards/{{class_id}}/{{item}}
		boardId = pathSegment[3]
		item = pathSegment[4]
		return
	} else if len(pathSegment) == 4 {
		// /{{version}}/boards/{{class_id}}
		boardId = pathSegment[3]
		return
	} else if len(pathSegment) == 3 {
		// /{{version}}/boards
		// Should not be reach...
		return
	}
	delivery.logger.Warningf("parseBoardPath got malform path: %v", path)
	return
}
