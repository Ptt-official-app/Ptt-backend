package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

// getUsers is a http handler function which will rewrite to correct route
func (delivery *Delivery) getUsers(w http.ResponseWriter, r *http.Request) {
	userID, item, err := parseUserPath(r.URL.Path)
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

	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadUserInformation},
		map[string]string{
			"user_id": userID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataMap, err := delivery.usecase.GetUserInformation(context.Background(), userID)
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
	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadUserInformation},
		map[string]string{
			"user_id": userID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataItems, err := delivery.usecase.GetUserFavorites(context.Background(), userID)
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
	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadUserInformation},
		map[string]string{
			"user_id": userID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// return need fix
	dataItems, err := delivery.usecase.GetUserArticles(context.Background(), userID)
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

	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadUserInformation},
		map[string]string{
			"user_id": userID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataMap, err := delivery.usecase.GetUserPreferences(context.Background(), userID)
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

	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadUserInformation},
		map[string]string{
			"user_id": userID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataItems, err := delivery.usecase.GetUserComments(context.Background(), userID)
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
