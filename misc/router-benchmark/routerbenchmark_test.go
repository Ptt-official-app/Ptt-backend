package router_benchmark

import (
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"strings"
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

// getBoards is the handler for `/v1/boards` with GET method
func getBoards(w http.ResponseWriter, r *http.Request) {
	boardID, item, filename, _ := parseBoardPath(r.URL.Path)
	if boardID == "" {
		//delivery.getBoardList(w, r)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// get single board
	if item == "information" {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if item == "settings" {
		//delivery.getBoardSettings(w, r, boardID)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if item == "articles" {
		if filename == "" {
			getBoardArticles(w, r, boardID)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		return
	} else if item == "treasures" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 404
	w.WriteHeader(http.StatusNotFound)
}

// parseBoardPath covert url path from /v1/boards/SYSOP/article to
// {SYSOP, article) or /v1/boards to {,}
func parseBoardPath(path string) (boardID string, item string, filename string, err error) {
	pathSegment := strings.Split(path, "/")

	if len(pathSegment) >= 6 {
		boardID = pathSegment[3]
		item = pathSegment[4]
		filename = pathSegment[5]
		return
	} else if len(pathSegment) == 5 {
		boardID = pathSegment[3]
		item = pathSegment[4]
		return
	} else if len(pathSegment) == 4 {
		boardID = pathSegment[3]
		return
	} else if len(pathSegment) == 3 {
		return
	}
	return
}

// getBoardArticles handles request with `/v1/boards/SYSOP/articles` and will return
// article list to client
func getBoardArticles(w http.ResponseWriter, _ *http.Request, boardID string) {
	_ = boardID
	w.WriteHeader(200)
}

// getBoards is the handler for `/v1/boards` with GET method
func getBoardArticles_httprouter(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	boardID := params.ByName("boardID")
	_ = boardID
	w.WriteHeader(200)
}

// getBoards is the handler for `/v1/boards` with GET method
func getBoardArticles_gorillamux(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	boardID := params["boardID"]
	_ = boardID
	w.WriteHeader(200)
}
