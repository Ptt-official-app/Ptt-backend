package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (delivery *Delivery) createBoard(w http.ResponseWriter, r *http.Request) {
	token := delivery.getTokenFromRequest(r)
	if err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionCreateBoard}, nil); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, writeErr := w.Write(NewPermissionError(r, err))
		if writeErr != nil {
			delivery.logger.Errorf("write board creation permission error: %w", writeErr)
		}
		return
	}

	boardID := r.PostFormValue("board_id")
	title := r.PostFormValue("title")
	board, err := delivery.usecase.CreateBoard(context.Background(), boardID, title)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidBoard) || errors.Is(err, repository.ErrBoardExists) {
			w.WriteHeader(http.StatusBadRequest)
			_, writeErr := w.Write(NewNewBoardError(r))
			if writeErr != nil {
				delivery.logger.Errorf("write new board validation error: %w", writeErr)
			}
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write(NewServerError(r, err))
		if writeErr != nil {
			delivery.logger.Errorf("write new board server error: %w", writeErr)
		}
		return
	}

	response := map[string]interface{}{
		"data": marshalBoardHeader(board),
	}
	body, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(body); err != nil {
		delivery.logger.Errorf("write new board response: %w", err)
	}
}
