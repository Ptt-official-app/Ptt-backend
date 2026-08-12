package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

type favoriteHTTPUsecase struct {
	MockUsecase
	err          error
	called       bool
	token        string
	userID       string
	favoriteType string
	boardID      string
	title        string
}

func (u *favoriteHTTPUsecase) AddUserFavorite(_ context.Context, token, userID, favoriteType, boardID, title string) ([]interface{}, error) {
	u.called = true
	u.token = token
	u.userID = userID
	u.favoriteType = favoriteType
	u.boardID = boardID
	u.title = title
	return []interface{}{}, u.err
}

func TestRouteAddUserFavorite(t *testing.T) {
	mockUsecase := &favoriteHTTPUsecase{}
	delivery := NewHTTPDelivery(mockUsecase)

	form := url.Values{}
	form.Set("action", "add_favorite")
	form.Set("type", "board")
	form.Set("board_id", "SYSOP")
	req := httptest.NewRequest(http.MethodPost, "/v1/users/pichu/favorites", strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "bearer test-token")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	delivery.routeUsers(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if !mockUsecase.called {
		t.Fatal("AddUserFavorite was not called")
	}
	if mockUsecase.token != "test-token" || mockUsecase.userID != "pichu" || mockUsecase.favoriteType != "board" || mockUsecase.boardID != "SYSOP" {
		t.Fatalf("unexpected args: token=%q user=%q type=%q board=%q", mockUsecase.token, mockUsecase.userID, mockUsecase.favoriteType, mockUsecase.boardID)
	}
}

func TestPostUserFavoritesStatusCodes(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "unauthorized", err: usecase.ErrFavoriteUnauthorized, want: http.StatusUnauthorized},
		{name: "forbidden", err: usecase.ErrFavoriteForbidden, want: http.StatusForbidden},
		{name: "user not found", err: usecase.ErrFavoriteUserNotFound, want: http.StatusNotFound},
		{name: "invalid favorite", err: usecase.ErrInvalidFavorite, want: http.StatusBadRequest},
		{name: "server error", err: errors.New("disk failure"), want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &favoriteHTTPUsecase{err: tt.err}
			delivery := NewHTTPDelivery(mockUsecase)
			form := url.Values{}
			form.Set("action", "add_favorite")
			form.Set("type", "line")
			req := httptest.NewRequest(http.MethodPost, "/v1/users/pichu/favorites", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			delivery.postUserFavorites(rr, req, "pichu")
			if rr.Code != tt.want {
				t.Fatalf("status = %d, want %d; body=%s", rr.Code, tt.want, rr.Body.String())
			}
		})
	}
}

func TestPostUserFavoritesRejectsUnknownAction(t *testing.T) {
	delivery := NewHTTPDelivery(&favoriteHTTPUsecase{})
	form := url.Values{}
	form.Set("action", "delete_favorite")
	req := httptest.NewRequest(http.MethodPost, "/v1/users/pichu/favorites", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	delivery.postUserFavorites(rr, req, "pichu")
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}
