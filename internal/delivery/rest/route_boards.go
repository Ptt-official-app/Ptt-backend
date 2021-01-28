package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PichuChen/go-bbs"
)

// routeBoards is the handler for `/v1/boards`
func (rest *restHandler) routeBoards(w http.ResponseWriter, r *http.Request) {
	rest.logger.Debugf("routeBoards: %v", r)

	// TODO: Check IP Flowspeed
	if r.Method == "GET" {
		rest.getBoards(w, r)
		return
	}

}

// getBoards is the handler for `/v1/boards` with GET method
func (rest *restHandler) getBoards(w http.ResponseWriter, r *http.Request) {
	rest.logger.Debugf("getBoards: %v", r)
	boardId, item, filename, err := rest.parseBoardPath(r.URL.Path)
	if boardId == "" {
		rest.getBoardList(w, r)
		return
	}
	// get single board
	if item == "information" {
		rest.getBoardInformation(w, r, boardId)
		return
	} else if item == "articles" {
		if filename == "" {
			rest.getBoardArticles(w, r, boardId)
		} else {
			rest.getBoardArticlesFile(w, r, boardId, filename)
		}
		return
	} else if item == "treasures" {
		rest.getBoardTreasures(w, r, boardId)
		return
	}

	// 404
	w.WriteHeader(http.StatusNotFound)

	rest.logger.Noticef("board id: %v not exist but be queried, info: %v err: %v", boardId, item, err)
}

func (rest *restHandler) getBoardList(w http.ResponseWriter, r *http.Request) {
	rest.logger.Debugf("getBoardList: %v", r)

	token := rest.getTokenFromRequest(r)
	userId, err := rest.getUserIdFromToken(token)
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
	for _, b := range rest.boardRepo.GetBoards(context.Background()) {
		// TODO: Show Board by user level
		if b.IsClass() {
			continue
		}
		if !shouldShowOnUserLevel(b, userId) {
			continue
		}
		jb, _ := json.Marshal(b)
		rest.logger.Debugf("marshal board: %v", string(jb))
		dataList = append(dataList, marshalBoardHeader(b))
	}

	responseMap := map[string]interface{}{
		"data": dataList,
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)

}

func (rest *restHandler) getBoardInformation(w http.ResponseWriter, r *http.Request, boardId string) {
	rest.logger.Debugf("getBoardInformation: %v", r)
	token := rest.getTokenFromRequest(r)
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

	brd, err := rest.findBoardHeaderById(boardId)
	if err != nil {
		// TODO: record error
		rest.logger.Warningf("find board %s failed: %v", boardId, err)
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
func (rest *restHandler) parseBoardPath(path string) (boardId string, item string, filename string, err error) {
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
	rest.logger.Warningf("parseBoardPath got malform path: %v", path)
	return

}

func (rest *restHandler) findBoardHeaderById(boardId string) (bbs.BoardRecord, error) {
	for _, it := range rest.boardRepo.GetBoards(context.Background()) {
		if boardId == it.BoardId() {
			return it, nil
		}
	}
	return nil, fmt.Errorf("board record not found")

}
