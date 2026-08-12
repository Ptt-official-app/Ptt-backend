package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (delivery *Delivery) postUserFavorites(w http.ResponseWriter, r *http.Request, userID string) {
	if r.PostFormValue("action") != "add_favorite" {
		writeFavoriteError(w, http.StatusBadRequest, "invalid_action", "action must be add_favorite")
		return
	}

	items, err := delivery.usecase.AddUserFavorite(
		context.Background(),
		delivery.getTokenFromRequest(r),
		userID,
		r.PostFormValue("type"),
		r.PostFormValue("board_id"),
		r.PostFormValue("title"),
	)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrFavoriteUnauthorized):
			writeFavoriteError(w, http.StatusUnauthorized, "permission_error", err.Error())
		case errors.Is(err, usecase.ErrFavoriteForbidden):
			writeFavoriteError(w, http.StatusForbidden, "permission_error", err.Error())
		case errors.Is(err, usecase.ErrFavoriteUserNotFound):
			writeFavoriteError(w, http.StatusNotFound, "user_not_found", err.Error())
		case errors.Is(err, usecase.ErrInvalidFavorite):
			writeFavoriteError(w, http.StatusBadRequest, "invalid_favorite", err.Error())
		default:
			delivery.logger.Errorf("add favorite for %s failed: %v", userID, err)
			writeFavoriteError(w, http.StatusInternalServerError, "server_error", err.Error())
		}
		return
	}

	response, err := json.MarshalIndent(map[string]interface{}{
		"data": map[string]interface{}{
			"items": items,
		},
	}, "", "  ")
	if err != nil {
		writeFavoriteError(w, http.StatusInternalServerError, "server_error", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		delivery.logger.Errorf("write add favorite response: %v", err)
	}
}

func writeFavoriteError(w http.ResponseWriter, status int, code, description string) {
	body, _ := json.MarshalIndent(map[string]string{
		"error":             code,
		"error_description": description,
	}, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
