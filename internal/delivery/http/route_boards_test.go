package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestGetBoardList (t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/boards", nil)
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
	r.HandleFunc("/v1/boards", delivery.routeBoards)
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
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
