package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (delivery *httpDelivery) getBoardList(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getBoardList: %v", r)

	token := delivery.getTokenFromRequest(r)
	userID, err := delivery.usecase.GetUserIDFromToken(token)
	if err != nil {
		// user permission error
		// Support Guest?
		if !supportGuest() {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"token_invalid"}`))
			return
		}
		userID = "guest" // TODO: use const variable
	}

	boards := delivery.usecase.GetBoards(context.Background(), userID)

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

func (delivery *httpDelivery) getBoardInformation(w http.ResponseWriter, r *http.Request, boardID string) {
	delivery.logger.Debugf("getBoardInformation: %v", r)
	token := delivery.getTokenFromRequest(r)
	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	brd, err := delivery.usecase.GetBoardByID(context.Background(), boardID)
	if err != nil {
		// TODO: record error
		delivery.logger.Warningf("find board %s failed: %v", boardID, err)
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_board_error",
			"error_description": "get board for " + boardID + " failed",
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
