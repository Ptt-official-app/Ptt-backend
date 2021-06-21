package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

// getUsers is a http handler function which will rewrite to correct route
func (delivery *Delivery) getUsers(w http.ResponseWriter, r *http.Request) {
	userID, item, itemID, err := parseUserPath(r.URL.Path)
	switch item {
	case "information":
		delivery.getUserInformation(w, r, userID)
	case "favorites":
		delivery.getUserFavorites(w, r, userID)
	case "articles":
		delivery.getUserArticles(w, r, userID)
	case "preferences":
		delivery.getUserPreferences(w, r, userID)
	case "comments":
		delivery.getUserComments(w, r, userID)
	case "drafts":
		delivery.getUserDrafts(w, r, userID, itemID)
	default:
		delivery.logger.Noticef("user id: %v not exist but be queried, info: %v err: %v", userID, item, err)
		w.WriteHeader(http.StatusNotFound)
	}
}

func (delivery *Delivery) postUsers(w http.ResponseWriter, r *http.Request) {
	userID, item, itemID, err := parseUserPath(r.URL.Path)
	switch item {
	case "drafts":
		delivery.postUserDrafts(w, r, userID, itemID)
	default:
		delivery.logger.Noticef("user id: %v not exist but be queried, info: %v err: %v", userID, item, err)
		w.WriteHeader(http.StatusNotFound)
	}
}

// getUserInformation is a http handler function which will writes the information including last visit time of user with userID
// to w. request path should be /v1/users/{{user_id}}/information
// Please see: https://pttapp.cc/swagger/#/%E4%BD%BF%E7%94%A8%E8%80%85%E9%83%A8%E5%88%86/get_v1_users__user_id__information
func (delivery *Delivery) getUserInformation(w http.ResponseWriter, r *http.Request, userID string) {
	token := delivery.getTokenFromRequest(r)

	ctx := context.Background()

	err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionReadUserInformation}, map[string]string{
		"user_id": userID,
	})

	if err != nil {
		delivery.logger.Warningf("unauthorized get user information for %s: %v", userID, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataMap, err := delivery.usecase.GetUserInformation(ctx, userID)
	if err != nil {
		delivery.logger.Warningf("get user information for %s failed: %v", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_userrec_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getUserInformation error response err: %w", err)
		}
		return
	}

	responseMap := map[string]interface{}{
		"data": dataMap,
	}
	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	_, err = w.Write(responseByte)
	if err != nil {
		delivery.logger.Errorf("getUserInformation success response err: %w", err)
	}
}

// getUserFavorites is a http handler function which will get favorite list of user with userID
// to w. request path should be /v1/users/{{user_id}}/favorites
// Please see: https://pttapp.cc/swagger/#/%E4%BD%BF%E7%94%A8%E8%80%85%E9%83%A8%E5%88%86/get_v1_users__user_id__favorites
func (delivery *Delivery) getUserFavorites(w http.ResponseWriter, r *http.Request, userID string) {
	token := delivery.getTokenFromRequest(r)
	ctx := context.Background()
	err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionReadUserInformation}, map[string]string{
		"user_id": userID,
	})

	if err != nil {
		m := map[string]string{
			"error":             "permission_error",
			"error_description": "no permission",
		}

		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getUserFavorites error response err: %w", err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataItems, err := delivery.usecase.GetUserFavorites(ctx, userID)
	if err != nil {
		delivery.logger.Errorf("failed to get user favorites: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": dataItems,
		},
	}

	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	_, err = w.Write(responseByte)
	if err != nil {
		delivery.logger.Errorf("getBoardInformation write success response err: %w", err)
	}
}

// getUserArticles is a http handler function which will get user's articles list of user with userID
// to w. request path should be /v1/users/{{user_id}}/articles
func (delivery *Delivery) getUserArticles(w http.ResponseWriter, r *http.Request, userID string) {
	token := delivery.getTokenFromRequest(r)
	ctx := context.Background()
	err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionReadUserInformation}, map[string]string{
		"user_id": userID,
	})

	if err != nil {
		m := map[string]string{
			"error":             "permission_error",
			"error_description": "no permission",
		}

		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getUserArticles error response err: %w", err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// return need fix
	dataItems, err := delivery.usecase.GetUserArticles(ctx, userID)
	if err != nil {
		delivery.logger.Errorf("failed to get user's articles: %s\n", err)
	}

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": dataItems,
		},
	}

	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	_, err = w.Write(responseByte)
	if err != nil {
		delivery.logger.Errorf("getUserFavorites success response err: %w", err)
	}
}

// getUserPreferences is a http handler function which will get user's preferences list of user with userID
// to w. request path should be /v1/users/{{user_id}}/preferences
func (delivery *Delivery) getUserPreferences(w http.ResponseWriter, r *http.Request, userID string) {
	token := delivery.getTokenFromRequest(r)
	ctx := context.Background()

	err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionReadUserInformation}, map[string]string{
		"user_id": userID,
	})

	if err != nil {
		m := map[string]string{
			"error":             "permission_error",
			"error_description": "no permission",
		}

		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getUserPreferences error response err: %w", err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataMap, err := delivery.usecase.GetUserPreferences(ctx, userID)
	if err != nil {
		// TODO: record error
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_userrec_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getUserPreferences error response err: %w", err)
		}
		return
	}

	responseMap := map[string]interface{}{
		"data": dataMap,
	}
	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	_, err = w.Write(responseByte)
	if err != nil {
		delivery.logger.Errorf("getUserPreferences success response err: %w", err)
	}
}

// getUserComments is a http handler function which will get history comments of user with userID
// to w. request path should be /v1/users/{{user_id}}/comments
func (delivery *Delivery) getUserComments(w http.ResponseWriter, r *http.Request, userID string) {
	token := delivery.getTokenFromRequest(r)
	ctx := context.Background()

	err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionReadUserInformation}, map[string]string{
		"user_id": userID,
	})

	if err != nil {
		m := map[string]string{
			"error":             "permission_error",
			"error_description": "no permission",
		}

		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getUserComments error response err: %w", err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataItems, err := delivery.usecase.GetUserComments(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_userrec_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getUserComments error response err: %w", err)
		}
		return
	}

	dataList := make([]interface{}, 0, len(dataItems))
	for _, board := range dataItems {
		dataList = append(dataList, map[string]interface{}{
			"CommentOrder": board.CommentOrder(),
			"CommentOwner": board.CommentOwner(),
			"COmmentTime":  board.CommentTime(),
		})
	}

	responseMap := map[string]interface{}{
		"data": dataList,
	}

	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")

	_, err = w.Write(responseByte)
	if err != nil {
		delivery.logger.Errorf("getUserFavorites success response err: %w", err)
	}
}

func (delivery *Delivery) getUserDrafts(w http.ResponseWriter, r *http.Request, userID string, draftID string) {
	token := delivery.getTokenFromRequest(r)
	ctx := context.Background()

	err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionReadUserInformation}, map[string]string{
		"user_id": userID,
	})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	buf, err := delivery.usecase.GetUserDrafts(ctx, userID, draftID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_userrec_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getUserDrafts error response err: %w", err)
		}
		return
	}

	bufStr := base64.StdEncoding.EncodeToString(buf.Raw())

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"raw": bufStr,
		},
	}

	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(responseByte)
	if err != nil {
		delivery.logger.Errorf("getUserDrafts success response err: %w", err)
	}
}

func (delivery *Delivery) postUserDrafts(w http.ResponseWriter, r *http.Request, userID string, draftID string) {
	delivery.logger.Debugf("postUserDrafts: %v", r)

	action := r.PostFormValue("action")
	switch action {
	case "update_draft":
		delivery.updateUserDraft(w, r, userID, draftID)
	case "delete_draft":
		delivery.deleteUserDraft(w, r, userID, draftID)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (delivery *Delivery) updateUserDraft(w http.ResponseWriter, r *http.Request, userID string, draftID string) {
	token := delivery.getTokenFromRequest(r)

	ctx := context.Background()
	err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionUpdateDraft}, map[string]string{
		"user_id": userID,
	})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	raw := r.PostFormValue("raw")
	if raw == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	text, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		delivery.logger.Errorf("failed to decode raw text: %s\n", err)
		return
	}

	buf, err := delivery.usecase.UpdateUserDraft(ctx, userID, draftID, text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "update_userrec_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("updateDraft error response err: %w", err)
		}
		return
	}

	bufStr := base64.StdEncoding.EncodeToString(buf.Raw())

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"raw": bufStr,
		},
	}

	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(responseByte)
	if err != nil {
		delivery.logger.Errorf("updateDraft success response err: %w", err)
	}
}

func (delivery *Delivery) deleteUserDraft(w http.ResponseWriter, r *http.Request, userID string, draftID string) {
	token := delivery.getTokenFromRequest(r)
	ctx := context.Background()

	err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionDeleteDraft}, map[string]string{
		"user_id": userID,
	})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = delivery.usecase.DeleteUserDraft(ctx, userID, draftID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "delete_userrec_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("deleteDraft error response err: %w", err)
		}
		return
	}

	responseMap := map[string]interface{}{
		"data": struct{}{},
	}

	responseByte, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(responseByte)
	if err != nil {
		delivery.logger.Errorf("deleteDraft success response err: %w", err)
	}
}
