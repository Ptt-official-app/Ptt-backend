package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchBoardsTreasures(t *testing.T) {
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)
	req, err := http.NewRequest("GET", "/v1/boards/1/treasures/D333/articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("v1/boards/", delivery.routeBoards)
	r.ServeHTTP(rr, req)
}
