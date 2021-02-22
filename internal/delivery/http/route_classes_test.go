package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetClassesList(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)
	req, err := http.NewRequest("GET", "/v1/classes/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/classes/", delivery.routeClasses)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responseMap := map[string][]interface{}{}
	json.Unmarshal(rr.Body.Bytes(), &responseMap)
	t.Logf("got response %v", rr.Body.String())
	responseData := responseMap["data"]
	for i, d := range responseData {
		board := d.(map[string]interface{})
		expectedID := fmt.Sprintf("%v", i+1)
		if board["type"] == "class" && board["id"] != expectedID {
			t.Errorf("handler returned unexpected body, id not match: got %v want %v",
				board, expectedID)
		}
	}
}
