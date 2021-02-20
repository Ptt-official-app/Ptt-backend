package http

import (
	"context"
	// "encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (delivery *httpDelivery) getBoardTreasures(w http.ResponseWriter, r *http.Request, boardID string) {
	delivery.logger.Debugf("getBoardTreasures: %v", r)
	token := delivery.getTokenFromRequest(r)
	_, treasuresID, filename, err := delivery.parseBoardTreasurePath(r.URL.Path)
	if err != nil {
		delivery.logger.Warningf("parseBoardTreasurePath error: %v", err)
		// TODO return 400?
	}
	if filename != "" {
		// get file
		delivery.getBoardTreasuresFile(w, r, boardID, treasuresID, filename)
		return
	}

	// Check permission for board
	err = delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadTreasureInformation},
		map[string]string{
			"board_id":    boardID,
			"treasure_id": strings.Join(treasuresID, ","),
		})
	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	treasures := delivery.usecase.GetBoardTreasures(context.Background(), boardID, treasuresID)
	delivery.logger.Debugf("fh: %v", treasures)

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": treasures,
		},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	_, _ = w.Write(b)

}

func (delivery *httpDelivery) getBoardTreasuresFile(w http.ResponseWriter, r *http.Request, boardID string, treasuresID []string, filename string) {
	delivery.logger.Debugf("getBoardTreasuresFile %v board: %v, treasuresID: %v, filename: %v", r, boardID, treasuresID, filename)

	w.WriteHeader(http.StatusNotImplemented)
}
