package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

// forwardArticle handles request with `/v1/boards/{boardID}/articles/{filename}` and will
// forward article to either email or board
func (delivery *Delivery) forwardArticle(w http.ResponseWriter, r *http.Request, boardID, filename string) {
	delivery.logger.Debugf("forwardArticle: %v", r)

	toEmail := r.PostFormValue("email")
	toBoard := r.PostFormValue("board")
	// either email or boardID must be valid
	if toEmail == "" && toBoard == "" {
		w.WriteHeader(500)
		return
	}

	token := delivery.getTokenFromRequest(r)

	// Check permission for forwarding article
	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionForwardArticle},
		map[string]string{
			"board_id":   boardID,
			"article_id": filename,
		})
	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := delivery.usecase.GetUserIDFromToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	destinations := []usecase.Forward{
		&usecase.ForwardToBoard{
			Board: toBoard,
		},
		&usecase.ForwardToEmail{
			Email: toEmail,
		},
	}

	for _, dest := range destinations {
		_, err = delivery.usecase.ForwardArticle(
			context.Background(),
			userID,
			boardID,
			filename,
			dest,
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("forwardArticle write success response err: %w", err)
	}
}
