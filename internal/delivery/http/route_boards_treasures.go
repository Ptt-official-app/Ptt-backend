package http

import (
	"context"
	// "encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

func (delivery *httpDelivery) getBoardTreasures(w http.ResponseWriter, r *http.Request, boardId string) {
	delivery.logger.Debugf("getBoardTreasures: %v", r)
	token := delivery.getTokenFromRequest(r)
	_, treasuresId, filename, err := delivery.parseBoardTreasurePath(r.URL.Path)
	if err != nil {
		delivery.logger.Warningf("parseBoardTreasurePath error: %v", err)
		// TODO return 400?
	}
	if filename != "" {
		// get file
		delivery.getBoardTreasuresFile(w, r, boardId, treasuresId, filename)
		return
	}

	// Check permission for board
	err = checkTokenPermission(token,
		[]permission{PermissionReadTreasureInformation},
		map[string]string{
			"board_id":    boardId,
			"treasure_id": strings.Join(treasuresId, ","),
		})
	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	treasures := delivery.boardUsecase.GetBoardTreasures(context.Background(), boardId, treasuresId)
	delivery.logger.Debugf("fh: %v", treasures)

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": treasures,
		},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)

}

func (delivery *httpDelivery) getBoardTreasuresFile(w http.ResponseWriter, r *http.Request, boardId string, treasuresId []string, filename string) {
	delivery.logger.Debugf("getBoardTreasuresFile %v board: %v, treasuresId: %v, filename: %v", r, boardId, treasuresId, filename)

	w.WriteHeader(http.StatusNotImplemented)
}

// parseBoardTreasurePath parse covert url path from /v1/boards/SYSOP/article to
// {SYSOP, article) or /v1/boards to {,}
func (delivery *httpDelivery) parseBoardTreasurePath(path string) (boardId string, treasuresId []string, filename string, err error) {
	pathSegment := strings.Split(path, "/")

	if len(pathSegment) == 6 {
		// /{{version}}/boards/{{board_id}}/treasures/articles
		boardId = pathSegment[3]
		treasuresId = []string{}
		filename = ""
		return
	} else if len(pathSegment) >= 7 {
		// /{{version}}/boards/{{board_id}}/treasures/{{treasures_id ... }}/articles
		// or
		// /{{version}}/boards/{{board_id}}/treasures/{{treasures_id ... }}/articles/{{filename}}
		boardId = pathSegment[3]
		if pathSegment[len(pathSegment)-1] == "articles" {
			treasuresId = pathSegment[5 : len(pathSegment)-1]
			filename = ""
		} else {
			treasuresId = pathSegment[5 : len(pathSegment)-2]
			filename = pathSegment[len(pathSegment)-1]
		}
		return
	}
	// should not be reached
	delivery.logger.Warningf("parseBoardTreasurePath got malform path: %v", path)
	return

}
