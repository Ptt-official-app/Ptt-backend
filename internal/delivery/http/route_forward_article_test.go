package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
)

// TestForwardArticleBadRequest test request post `/v1/boards/{}/articles/{}` post
// with no body
func TestForwardArticleBadRequest(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase, logging.DefaultDummyLogger)

	req, err := http.NewRequest("POST", "/v1/boards/test/articles/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/", delivery.routeBoards)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestForwardArticleBadRequest test request post `/v1/boards/{}/articles/{}` post
// without valid email and board
func TestForwardArticleInternalServerError(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase, logging.DefaultDummyLogger)

	v := url.Values{}
	v.Set("action", "forward_article")
	v.Set("board", "")
	v.Set("email", "")
	t.Logf("testing body: %v", v)
	req, err := http.NewRequest("POST", "/v1/boards/test/articles/test", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/", delivery.routeBoards)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != 500 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 500)
	}
}

func TestForwardArticleResponse(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase, logging.DefaultDummyLogger)

	v := url.Values{}
	v.Set("action", "forward_article")
	v.Set("board", "SYSOP")
	v.Set("email", "pichu@tih.tw")
	t.Logf("testing body: %v", v)
	req, err := http.NewRequest("POST", "/v1/boards/test/articles/test", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/", delivery.routeBoards)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
