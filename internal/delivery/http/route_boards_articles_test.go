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
	json.Unmarshal(rr.Body.Bytes(), &actualResponseMap)
	t.Logf("got response: %v", rr.Body.String())

	actualResponseDataList := actualResponseMap["data"].(map[string]interface{})
	actualResponseItems := actualResponseDataList["items"].([]interface{})

	actualResponseData := actualResponseItems[0].(map[string]interface{})
	if _, ok := actualResponseData["filename"]; !ok {
		t.Error("expect response has index \"filename\"")
	}
	if _, ok := actualResponseData["modified_time"]; !ok {
		t.Error("expect response has index \"modified_time\"")
	}
	if _, ok := actualResponseData["recommend_count"]; !ok {
		t.Error("expect response has index \"recommend_count\"")
	}
	if _, ok := actualResponseData["post_date"]; !ok {
		t.Error("expect response has index \"post_date\"")
	}
	if _, ok := actualResponseData["title"]; !ok {
		t.Error("expect response has index \"title\"")
	}
	if _, ok := actualResponseData["money"]; !ok {
		t.Error("expect response has index \"money\"")
	}
	if _, ok := actualResponseData["owner"]; !ok {
		t.Error("expect response has index \"owner\"")
	}
	if _, ok := actualResponseData["url"]; !ok {
		t.Error("expect response has index \"url\"")
	}
}
