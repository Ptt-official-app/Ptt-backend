package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestGetUserInformation is a test function which will test getUserInformation (/v1/users/{{user_id}}/favorites)
// Please see: https://pttapp.cc/swagger/#/%E4%BD%BF%E7%94%A8%E8%80%85%E9%83%A8%E5%88%86/get_v1_users__user_id__information
func TestGetUserInformation(t *testing.T) {
	userID := "id"
	expectedData := map[string]interface{}{
		"user_id":              userID,
		"nickname":             "",
		"realname":             "",
		"number_of_login_days": "0",
		"number_of_posts":      "0",
		"number_of_badposts":   "0",
		"money":                "0",
		"money_description":    "債台高築",
		"last_login_time":      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		"last_login_ipv4":      "127.0.0.1",
		"last_login_ip":        "127.0.0.1",
		"last_login_country":   "",
		"mailbox_description":  "",
		"chess_status":         map[string]interface{}{},
		"plan":                 map[string]interface{}{},
	}
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)
	req, err := http.NewRequest("GET", "/v1/users/SYSOP/information", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/users/", delivery.routeUsers)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responsedMap := map[string]interface{}{}
	err = json.Unmarshal(rr.Body.Bytes(), &responsedMap)
	if err != nil {
		t.Errorf("get unexpect json: %w", err)
	}

	t.Logf("got response %v", rr.Body.String())
	responsedData := responsedMap["data"].(map[string]interface{})

	responseB, _ := json.MarshalIndent(responsedData, "", "  ")
	expectedB, _ := json.MarshalIndent(expectedData, "", "  ")
	if string(responseB) != string(expectedB) {
		t.Errorf("response not match: got %v want %v", responseB, expectedB)
	}
}

// TestParseUserPath is a test function which will test getUsers route mapping
func TestParseUserPath(t *testing.T) {

	type TestCase struct {
		input         string
		expectdUserID string
		expectdItem   string
	}

	cases := []TestCase{
		{
			input:         "/v1/users/Pichu/information",
			expectdUserID: "Pichu",
			expectdItem:   "information",
		},
		{
			input:         "/v1/users/Pichu/",
			expectdUserID: "Pichu",
			expectdItem:   "",
		},
		{
			input:         "/v1/users/Pichu",
			expectdUserID: "Pichu",
			expectdItem:   "",
		},
	}

	for index, c := range cases {
		input := c.input
		expectdUserID := c.expectdUserID
		expectdItem := c.expectdItem
		actualUserID, actualItem, err := parseUserPath(input)
		if err != nil {
			t.Errorf("error on index %d, got: %v", index, err)

		}

		if actualUserID != expectdUserID {
			t.Errorf("userId not match on index %d, expected: %v, got: %v", index, expectdUserID, actualUserID)
		}

		if actualItem != expectdItem {
			t.Errorf("item not match on index %d, expected: %v, got: %v", index, expectdItem, actualItem)
		}

	}

}

// TestGetUserInformation is a test function which will test getUserInformation (/v1/users/{{user_id}}/favorites)
// Please see: https://pttapp.cc/swagger/#/%E4%BD%BF%E7%94%A8%E8%80%85%E9%83%A8%E5%88%86/get_v1_users__user_id__favorites
func TestGetUserFavorite(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)
	req, err := http.NewRequest("GET", "/v1/users/SYSOP/favorites", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/users/", delivery.routeUsers)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responsedMap := map[string]interface{}{}
	err = json.Unmarshal(rr.Body.Bytes(), &responsedMap)
	if err != nil {
		t.Errorf("get unexpect json: %w", err)
	}

	t.Logf("got response %v", rr.Body.String())
	responsedData := responsedMap["data"].(map[string]interface{})
	items := responsedData["items"].([]interface{})
	firstItem := items[0].(map[string]interface{})

	expectBoardID := "test_board_001"
	if firstItem["board_id"].(string) != "test_board_001" {
		t.Errorf("handler returned unexpected body, board_id not match: got %v want board_id %v",
			firstItem["board_id"], expectBoardID)
	}
}

// TestGetUserPreference is a test function which will test getUserPreferences (/v1/users/{{user_id}}/preferences)
func TestGetUserPreference(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)
	req, err := http.NewRequest("GET", "/v1/users/SYSOP/preferences", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer"+token)

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/users/", delivery.routeUsers)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	responsedMap := map[string]interface{}{}
	err = json.Unmarshal(rr.Body.Bytes(), &responsedMap)
	if err != nil {
		t.Errorf("get unexpected json: %w", err)
	}

	t.Logf("got response %v", rr.Body.String())
	responsedData := responsedMap["data"].(map[string]interface{})
	firstItem := responsedData["favorite_no_highlight"]

	expectedValue := "false"
	if firstItem != expectedValue {
		t.Errorf("handler returned unexpected body, favorite_no_highlight not match: got %v want value %v",
			firstItem, expectedValue)
	}

}

// TestGetUserArticles is a test function which will test getUserArticles (/v1/users/{{user_id}}/articles)
func TestGetUserArticles(t *testing.T) {

	userID := "id"
	mockUsecase := NewMockUsecase()
	mockDelivery := NewHTTPDelivery(mockUsecase)

	req, err := http.NewRequest("GET", "/v1/users/user/articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := mockUsecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/users/", mockDelivery.routeUsers)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responsedMap := map[string]interface{}{}
	err = json.Unmarshal(rr.Body.Bytes(), &responsedMap)
	if err != nil {
		t.Errorf("get unexpect json: %w", err)
	}
	t.Logf("got response %v", rr.Body.String())
	if responsedMap["data"] == nil {
		t.Errorf("handler returned unexpected body, got %v want not nil",
			rr.Body.String())
	}
}

// TestGetUserComments is a test function which will test getUserComments (/v1/users/{{user_id}}/comments)
func TestGetUserComments(t *testing.T) {
	userID := "id"
	mockUsecase := NewMockUsecase()
	mockDelivery := NewHTTPDelivery(mockUsecase)

	req, err := http.NewRequest("GET", "/v1/users/user/comments", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := mockUsecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/users/", mockDelivery.routeUsers)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responsedMap := map[string]interface{}{}
	err = json.Unmarshal(rr.Body.Bytes(), &responsedMap)
	if err != nil {
		t.Errorf("get unexpect json: %w", err)
	}

	t.Logf("got response %v", rr.Body.String())
	responsedData := responsedMap["data"].(map[string]interface{})
	items := responsedData["items"].([]interface{})
	firstItem := items[0].(map[string]interface{})

	expectedValue := "SYSOP"
	if firstItem["board_id"].(string) != expectedValue {
		t.Errorf("handler returned unexpected body, favorite_no_highlight not match: got %v want value %v",
			firstItem, expectedValue)
	}
}
