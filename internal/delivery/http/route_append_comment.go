package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
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
			"raw": r.PostForm.Encode(),
			"parsed": map[string]interface{}{
				"is_header_modied": false,
				"author_id":        nil,
				"author_name":      nil,
				"title":            nil,
				"post_time":        nil,
				"board_name":       "", // todo: go-bbs articles 需實作新介面取得資訊
				"text": map[string]string{
					"text": "", // todo: // todo: go-bbs articles 需實作新介面取得資訊
				},
				"signature":    map[string]string{},
				"sender_info":  map[string]string{}, // todo: go-bbs articles 需實作新介面取得資訊(user info)
				"edit_records": map[string]string{},
				"push_records": []repository.PushRecord{res}, // TODO: support bbs.PushRecord instead
			},
		},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("getBoardTreasures write success response err: %w", err)
	}
}
