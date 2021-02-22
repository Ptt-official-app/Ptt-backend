package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"strconv"
)

func TestGetPopularBoardList(t *testing.T) {

	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)
	req, err := http.NewRequest("GET", "/v1/popular-boards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/popular-boards", delivery.routePopularBoards)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responseMap := map[string]interface{}{}
	json.Unmarshal(rr.Body.Bytes(), &responseMap)
	t.Logf("got response %v", rr.Body.String())
	responseData := responseMap["data"]
	popularBoards := responseData.(map[string]interface{})["items"].([]interface{})

	var prevNum int
	for i := range popularBoards {
		curr := popularBoards[i].(map[string]interface{})["number_of_user"]
		currNum, err := strconv.Atoi(curr.(string))
		if err != nil {
			t.Fatalf("handler returned unexpected body, invalid number_of_user: got %v",
				currNum)
		}
		
		if i > 0 && prevNum < currNum {
			t.Fatalf("handler returned unexpected body, invalid order: got %v before %v",
				prevNum, currNum)
		}
		prevNum = currNum
	}
}
