package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetBoardArticlesBadRequest(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	title := ""
	author := ""
	recommendCountGe := "qwerty"
	v := url.Values{}
	v.Set("title", title)
	v.Set("author", author)
	v.Set("recommend_count_ge", recommendCountGe)
	uri := fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req, err := http.NewRequest("GET", uri, nil)
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

func TestGetBoardArticlesResponse(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	titleContain := "test_posts"
	author := "test01"
	v := url.Values{}
	v.Set("title_contain", titleContain)
	v.Set("author", author)
	uri := fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req, err := http.NewRequest("GET", uri, nil)
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

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	actualResponseMap := map[string]interface{}{}
	if err := json.Unmarshal(rr.Body.Bytes(), &actualResponseMap); err != nil {
		t.Error(err.Error())
	}
}
