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

type deleteArticleHTTPUsecase struct {
	MockUsecase
	err      error
	called   bool
	token    string
	boardID  string
	filename string
}

func (u *deleteArticleHTTPUsecase) DeleteArticle(_ context.Context, token, boardID, filename string) error {
	u.called = true
	u.token = token
	u.boardID = boardID
	u.filename = filename
	return u.err
}

func TestRouteDeleteArticle(t *testing.T) {
	mockUsecase := &deleteArticleHTTPUsecase{}
	delivery := NewHTTPDelivery(mockUsecase)

	form := url.Values{}
	form.Set("action", "delete")
	req := httptest.NewRequest(http.MethodPost, "/v1/boards/test/articles/M.1.A.001", strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "bearer test-token")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	delivery.routeBoards(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if !mockUsecase.called {
		t.Fatal("DeleteArticle was not called")
	}
	if mockUsecase.token != "test-token" || mockUsecase.boardID != "test" || mockUsecase.filename != "M.1.A.001" {
		t.Fatalf("unexpected DeleteArticle args: token=%q board=%q filename=%q", mockUsecase.token, mockUsecase.boardID, mockUsecase.filename)
	}
}

func TestDeleteArticleStatusCodes(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "unauthorized", err: usecase.ErrDeleteArticleUnauthorized, want: http.StatusUnauthorized},
		{name: "forbidden", err: usecase.ErrDeleteArticleForbidden, want: http.StatusForbidden},
		{name: "not found", err: usecase.ErrArticleNotFound, want: http.StatusNotFound},
		{name: "server error", err: errors.New("disk failure"), want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &deleteArticleHTTPUsecase{err: tt.err}
			delivery := NewHTTPDelivery(mockUsecase)
			req := httptest.NewRequest(http.MethodPost, "/v1/boards/test/articles/M.1.A.001", nil)
			rr := httptest.NewRecorder()

			delivery.deleteArticle(rr, req, "test", "M.1.A.001")
			if rr.Code != tt.want {
				t.Fatalf("status = %d, want %d; body=%s", rr.Code, tt.want, rr.Body.String())
			}
		})
	}
}
