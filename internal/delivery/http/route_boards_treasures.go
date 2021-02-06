package http

import (
	"context"
	// "encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
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
	err = delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadTreasureInformation},
		map[string]string{
			"board_id":    boardId,
			"treasure_id": strings.Join(treasuresId, ","),
		})
	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	treasures := delivery.usecase.GetBoardTreasures(context.Background(), boardId, treasuresId)
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
