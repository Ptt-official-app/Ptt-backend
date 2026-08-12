package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (delivery *Delivery) deleteArticle(w http.ResponseWriter, r *http.Request, boardID, filename string) {
	token := delivery.getTokenFromRequest(r)
	err := delivery.usecase.DeleteArticle(context.Background(), token, boardID, filename)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrDeleteArticleUnauthorized):
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write(NewPermissionError(r, err))
		case errors.Is(err, usecase.ErrDeleteArticleForbidden):
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write(NewPermissionError(r, err))
		case errors.Is(err, usecase.ErrArticleNotFound):
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write(NewPathNotFoundError(r))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(NewServerError(r, err))
		}
		return
	}

	body, err := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"board_id": boardID,
			"filename": filename,
			"deleted":  true,
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		delivery.logger.Errorf("write delete article response: %w", err)
	}
}
