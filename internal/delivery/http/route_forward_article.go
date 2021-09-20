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
	toBoard := r.PostFormValue("board_id")
	// either email or boardID must be valid
	if toEmail == "" && toBoard == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(NewNoRequiredParameterError(r, "email or board_id"))
		if err != nil {
			delivery.logger.Errorf("write NewNoRequiredParameterError error: %w", err)
		}
		return
	}

	token := delivery.getTokenFromRequest(r)

	ctx := context.Background()

	userID, err := delivery.usecase.GetUserIDFromToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write(TokenPermissionError(r, err))
		if err != nil {
			delivery.logger.Errorf("write TokenPermissionError error: %w", err)
		}
		return
	}

	if toBoard != "" {
		// Check permission for whether article is allow forwarding `to` board
		err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionForwardArticleToBoard}, map[string]string{
			"board_id":   boardID,
			"to_board":   toBoard,
			"article_id": filename,
		})
		if err != nil {
			// TODO: record unauthorized access
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write(NewNoPermissionForCreateBoardArticlesError(r, toBoard))
			if err != nil {
				delivery.logger.Errorf("write NoPermissionForCreateBoardArticlesError error: %w", err)
			}
			return
		}

		_, err = delivery.usecase.ForwardArticleToBoard(
			ctx,
			userID,
			boardID,
			filename,
			toBoard,
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(NewNoPermissionForCreateBoardArticlesError(r, toBoard))
			if err != nil {
				delivery.logger.Errorf("write NoPermissionForCreateBoardArticlesError error: %w", err)
			}
			delivery.logger.Noticef("user %s forward file %s from board %s to board %s failed with error %s", userID, filename, boardID, toBoard, err.Error())
			return
		}
	}

	if toEmail != "" {
		err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionForwardArticleToEmail}, map[string]string{
			"board_id":   boardID,
			"to_email":   toEmail,
			"article_id": filename,
		})
		if err != nil {
			// TODO: record unauthorized access
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write(NewNoPermissionForForwardArticleToEmailError(r, filename, toEmail))
			if err != nil {
				delivery.logger.Errorf("write NoPermissionForForwardArticleToEmailError error: %w", err)
			}
			return
		}
		err = delivery.usecase.ForwardArticleToEmail(
			ctx,
			userID,
			boardID,
			filename,
			toEmail,
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write(NewNoPermissionForForwardArticleToEmailError(r, filename, toEmail))
			if err != nil {
				delivery.logger.Errorf("write NoPermissionForForwardArticleToEmailError error: %w", err)
			}
			delivery.logger.Noticef("user %s forward file %s from board %s to email %s failed with error %s", userID, filename, toEmail, toBoard, err.Error())
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
