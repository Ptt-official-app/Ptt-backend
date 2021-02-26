package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
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
	err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
	if err != nil {
		t.Errorf("get unexpect json: %w", err)
	}
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

func TestGetBoardList(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	req, err := http.NewRequest("GET", "/v1/boards/", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	w := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/", delivery.routeBoards)
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	actualResponseMap := map[string]interface{}{}
	err = json.Unmarshal(w.Body.Bytes(), &actualResponseMap)
	if err != nil {
		t.Errorf("get unexpect json: %w", err)
	}
	t.Logf("got response %v", w.Body.String())
	actualResponseDataList := actualResponseMap["data"].([]interface{})
	actualResponseData := actualResponseDataList[0].(map[string]interface{})

	if !strings.EqualFold(actualResponseData["title"].(string), "發電 ◎雜七雜八的垃圾") {
		t.Errorf("Title got %s, but excepted %s",
			actualResponseData["title"], "發電 ◎雜七雜八的垃圾")
	}
	if actualResponseData["number_of_user"] != "0" {
		t.Errorf("Number of users got %s, but excepted %s",
			actualResponseData["number_of_user"], "0")
	}
	if len(actualResponseData["moderators"].([]interface{})) != 0 {
		t.Errorf("Number of users got %s, but excepted %d",
			actualResponseData["moderators"], 0)
	}
	if !strings.EqualFold(actualResponseData["type"].(string), "board") {
		t.Errorf("Type got %s, but excepted %s",
			actualResponseData["type"], "board")
	}
}

func TestGetBoardInformation(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	req, err := http.NewRequest("GET", "/v1/boards/SYSOP/information", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	w := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/SYSOP/information", delivery.routeBoards)
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	actualResponseMap := map[string]interface{}{}
	err = json.Unmarshal(w.Body.Bytes(), &actualResponseMap)
	if err != nil {
		t.Errorf("get unexpect json: %w", err)
	}
	t.Logf("got response %v", w.Body.String())
	actualResponseData := actualResponseMap["data"].(map[string]interface{})

	if !strings.EqualFold(actualResponseData["title"].(string), "嘰哩 ◎站長好!") {
		t.Errorf("Title got %s, but excepted %s",
			actualResponseData["title"], "嘰哩 ◎站長好!")
	}
	if actualResponseData["number_of_user"] != "0" {
		t.Errorf("Number of users got %s, but excepted %s",
			actualResponseData["number_of_user"], "0")
	}
	if len(actualResponseData["moderators"].([]interface{})) != 0 {
		t.Errorf("Number of users got %s, but excepted %d",
			actualResponseData["moderators"], 0)
	}
	if !strings.EqualFold(actualResponseData["type"].(string), "board") {
		t.Errorf("Type got %s, but excepted %s",
			actualResponseData["type"], "board")
	}
}
