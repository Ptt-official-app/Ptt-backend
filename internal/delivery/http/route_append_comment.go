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

	updateUsefulness := false
	if appendType == "↑" || appendType == "↓" {
		updateUsefulness = true
	}
	if appendType == "" || (!updateUsefulness && text == "") {
		w.WriteHeader(500)
		return
	}

	token := delivery.getTokenFromRequest(r)

	userID, err := delivery.usecase.GetUserIDFromToken(token)
	if err != nil {
		delivery.logger.Warningf("GetUserIDFromToken for %s: %v", userID, err)
		w.WriteHeader(http.StatusUnauthorized)
		m := map[string]string{
			"error":             "get_user_id_from_token_error",
			"error_description": "get_user_id_from_token_error",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("appendComment error response err: %w", err)
		}
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
		delivery.logger.Warningf("unauthorized get user information for %s: %v", userID, err)
		w.WriteHeader(http.StatusUnauthorized)
		m := map[string]string{
			"error":             "permission_error",
			"error_description": "no permission",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("appendComment error response err: %w", err)
		}
		return
	}

	var res repository.PushRecord
	if updateUsefulness {
		res, err = delivery.usecase.UpdateUsefulness(
			ctx,
			userID,
			boardID,
			filename,
			appendType,
		)
	} else {
		res, err = delivery.usecase.AppendComment(
			ctx,
			userID,
			boardID,
			filename,
			appendType,
			text,
		)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m := map[string]string{
			"error":             "append_comment_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("AppendComment error response err: %w", err)
		}
		return
	}

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"raw": res.Text(),
			"parsed": map[string]interface{}{
				"is_header_modied": false,
				"author_id":        userID,
				"author_name":      userID,
				"title":            nil,
				"post_time":        res.Time().Format("2006-01-02 15:04:05"),
				"board_name":       boardID, // todo: go-bbs articles 需實作新介面取得資訊
				"text": map[string]string{
					"text": text, // todo: // todo: go-bbs articles 需實作新介面取得資訊
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
