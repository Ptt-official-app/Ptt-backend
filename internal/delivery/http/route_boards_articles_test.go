package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBoardArticlesBadRequest(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	title := ""
	author := ""
	recommendCountGe := "qwerty"
	url := fmt.Sprintf("/v1/boards/test/articles?title=%s&author=%s&recommend_count_ge=%s", title, author, recommendCountGe)
	req, err := http.NewRequest("GET", url, nil)
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
