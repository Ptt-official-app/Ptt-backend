package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (delivery *Delivery) publishPost(w http.ResponseWriter, r *http.Request, boardID string) {
	title := r.PostFormValue("title")
	article := r.PostFormValue("article")

	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(NewNoRequiredParameterError(r, "title"))
		if err != nil {
			delivery.logger.Errorf("write NewNoRequiredParameterError error: %w", err)
		}
		return
	}

	token := delivery.getTokenFromRequest(r)

	userID, err := delivery.usecase.GetUserIDFromToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write(TokenPermissionError(r, err))
		if err != nil {
			delivery.logger.Errorf("write NewNoRequiredParameterError error: %w", err)
		}
		return
	}

	ctx := context.Background()

	err = delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionCreateArticle}, map[string]string{
		"board_id": boardID,
		"user_id":  userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err2 := w.Write(NewNoPermissionForCreateBoardArticlesError(r, boardID))
		if err2 != nil {
			delivery.logger.Errorf("write no permission article error: %w", err)
		}
		return
	}

	record, err := delivery.usecase.CreateArticle(ctx, userID, boardID, title, article)
	// 改成 _ 避免 declared but not used
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err2 := w.Write(NewServerError(r, fmt.Errorf("create article error: %w", err)))
		if err2 != nil {
			delivery.logger.Errorf("write NewServerError error: %w", err)
		}
		return
	}

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"raw": r.PostForm.Encode(),
			"parsed": map[string]interface{}{
				"is_header_modied": false,
				"author_id":        userID,
				"board_name":       boardID, // todo: go-bbs articles 需實作新介面取得資訊
				"author_name":      record.Owner(),
				"title":            record.Title(),
				"post_time":        record.Date(),
				"text": map[string]string{
					"text": article, // todo: go-bbs articles 需實作新介面取得資訊
				},
				"signature":   map[string]string{},
				"sender_info": map[string]string{}, // todo: go-bbs articles 需實作新介面取得資訊(user info)
			},
		},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("getBoardArticles write success response err: %w", err)
	}
}
