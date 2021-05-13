package router_benchmark

import (
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Benchmark_ServeMux(b *testing.B) {
	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/", getBoards)
	runTest(b, r, rr)
}
func Benchmark_gorillamux(b *testing.B) {
	rr := httptest.NewRecorder()
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/v1/boards/{boardID}/articles", getBoardArticles_gorillamux)
	runTest(b, r, rr)
}
func Benchmark_httprouter(b *testing.B) {
	rr := httptest.NewRecorder()
	r := httprouter.New()
	r.HandlerFunc(http.MethodGet, "/v1/boards/:boardID/articles", getBoardArticles_httprouter)
	runTest(b, r, rr)
}

func runTest(b *testing.B, r http.Handler, rr *httptest.ResponseRecorder) {
	req, err := http.NewRequest("GET", "/v1/boards/SYSOP/articles", nil)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		r.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			b.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusOK)
		}
	}
}
