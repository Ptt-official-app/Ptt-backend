package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	TestGetUserInformation is a test function which will test getUserInformation (/v1/users/{{user_id}}/favorites)
	Please see: https://pttapp.cc/swagger/#/%E4%BD%BF%E7%94%A8%E8%80%85%E9%83%A8%E5%88%86/get_v1_users__user_id__favorites
*/
func TestGetUserInformation(t *testing.T) {

	userID := "id"
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
	if responsedData["user_id"] != userID {
		t.Errorf("handler returned unexpected body, user_id not match: got %v want userId %v",
			rr.Body.String(), userID)
	}
}

/*
	TestParseUserPath is a test function which will test getUsers route mapping
*/
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
