package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestAppendCommentBadRequest test request post `/v1/boards/{}/articles/{}` post
// with no body
func TestAppendCommentBadRequest(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

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

func TestAppendCommentResponse(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	v := url.Values{}
	v.Set("action", "append_comment")
	v.Set("type", "like")
	v.Set("text", "頭香")
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
