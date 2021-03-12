package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchBoardsTreasures(t *testing.T) {
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)
	req, err := http.NewRequest(http.MethodGet, "/v1/boards/1/treasures", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/", delivery.routeBoards)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responseMap := map[string]interface{}{}
	if err := json.Unmarshal(rr.Body.Bytes(), &responseMap); err != nil{
		t.Errorf("response body unmarshal error: %v", err)
	}
	responseData := responseMap["data"].(map[string]interface{})
	responseItems := responseData["items"].([]interface{})

	for _, treasures := range responseItems {
		treasure := treasures.(map[string]interface{})
		if _, ok := treasure["filename"]; !ok {
			t.Errorf("returned body filename not found.")
		}
		if _, ok := treasure["post_date"]; !ok {
			t.Errorf("returned body post_date not found.")
		}
		if _, ok := treasure["title"]; !ok {
			t.Errorf("returned body title not found.")
		}
		if _, ok := treasure["owner"]; !ok {
			t.Errorf("returned body owner not found")
		}
		if _, ok := treasure["url"]; !ok {
			t.Errorf("returned body url not found")
		}
	}


}
