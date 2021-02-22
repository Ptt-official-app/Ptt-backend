package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	json.Unmarshal(rr.Body.Bytes(), &responsedMap)
	t.Logf("got response %v", rr.Body.String())
	responsedData := responsedMap["data"].(map[string]interface{})
	if responsedData["user_id"] != userID {
		t.Errorf("handler returned unexpected body, user_id not match: got %v want userId %v",
			rr.Body.String(), userID)
	}
}

func TestParseUserPath(t *testing.T) {

	type TestCase struct {
		input         string
		expectdUserId string
		expectdItem   string
	}

	cases := []TestCase{
		{
			input:         "/v1/users/Pichu/information",
			expectdUserId: "Pichu",
			expectdItem:   "information",
		},
		{
			input:         "/v1/users/Pichu/",
			expectdUserId: "Pichu",
			expectdItem:   "",
		},
		{
			input:         "/v1/users/Pichu",
			expectdUserId: "Pichu",
			expectdItem:   "",
		},
	}

	for index, c := range cases {
		input := c.input
		expectdUserId := c.expectdUserId
		expectdItem := c.expectdItem
		actualUserId, actualItem, err := parseUserPath(input)
		if err != nil {
			t.Errorf("error on index %d, got: %v", index, err)

		}

		if actualUserId != expectdUserId {
			t.Errorf("userId not match on index %d, expected: %v, got: %v", index, expectdUserId, actualUserId)
		}

		if actualItem != expectdItem {
			t.Errorf("item not match on index %d, expected: %v, got: %v", index, expectdItem, actualItem)
		}

	}

}
