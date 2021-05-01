package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
)

func TestGetPopularArticles(t *testing.T) {
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase, logging.DefaultDummyLogger)
	req, err := http.NewRequest("GET", "/v1/popular-articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/popular-articles", delivery.routePopularArticles)
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
}
