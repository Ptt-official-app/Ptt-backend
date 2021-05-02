package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

// appendComment handles request with `/v1/boards/{boardID}/articles/{filename}` and will
// add comment to the article
func (delivery *Delivery) appendComment(w http.ResponseWriter, r *http.Request, boardID, filename string) {
	delivery.logger.Debugf("appendComment: %v", r)

	appendType := r.PostFormValue("type")
	text := r.PostFormValue("text")
	if appendType == "" || text == "" {
		w.WriteHeader(500)
		return
	}

	token := delivery.getTokenFromRequest(r)

	userID, err := delivery.usecase.GetUserIDFromToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx := context.Background()

	// Check permission for append comment
	err = delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionAppendComment}, map[string]string{
		"board_id":   boardID,
		"article_id": filename,
		"user_id":    userID,
	})
	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	res, err := delivery.usecase.AppendComment(
		ctx,
		userID,
		boardID,
		filename,
		appendType,
		text,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"raw":    r.PostForm.Encode(),
			"parsed": res,
		},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("getBoardTreasures write success response err: %w", err)
	}
}
