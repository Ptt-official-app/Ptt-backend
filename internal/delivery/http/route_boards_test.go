package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type exceptedMarshalBoardHeaderOfClassType struct {
	title string
	numberOfUser string
	lenOfModerators int
	boardType string
}

func TestGetBoardList (t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/boards/", nil)
	if err != nil {
		t.Fatal(err)
	}

	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

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
	json.Unmarshal(w.Body.Bytes(), &actualResponseMap)
	t.Logf("got response %v", w.Body.String())
	actualResponseDataList := actualResponseMap["data"].([]interface{})
	actualResponseData := actualResponseDataList[0].(map[string]interface{})

	exceptedBoard := &exceptedMarshalBoardHeaderOfClassType{
		title: "",
		numberOfUser: "0",
		lenOfModerators: 0,
		boardType: "class",
	}

	shouldBeEqualTo(actualResponseData, exceptedBoard, t)
}

func shouldBeEqualTo(actualResponseData map[string]interface{}, exceptedBoard *exceptedMarshalBoardHeaderOfClassType, t *testing.T) {
	if actualResponseData["title"] != exceptedBoard.title {
		t.Errorf("Title got %s, but excepted %s",
			actualResponseData["title"], "")
	}
	if actualResponseData["number_of_user"] != exceptedBoard.numberOfUser {
		t.Errorf("Number of users got %s, but excepted %s",
			actualResponseData["number_of_user"], "0")
	}
	if len(actualResponseData["moderators"].([]interface{})) != exceptedBoard.lenOfModerators {
		t.Errorf("Number of users got %s, but excepted %d",
			actualResponseData["moderators"], 0)
	}
	if actualResponseData["type"] != exceptedBoard.boardType {
		t.Errorf("Type got %s, but excepted %s",
			actualResponseData["type"], "class")
	}
}

func TestGetBoardInformation(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/boards/class/information", nil)
	if err != nil {
		t.Fatal(err)
	}

	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	w := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/class/information", delivery.routeBoards)
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
