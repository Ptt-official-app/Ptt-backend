package router_benchmark

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// getBoards is the handler for `/v1/boards` with GET method
func getBoardArticles_httprouter(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	boardID := params.ByName("boardID")
	_ = boardID
	w.WriteHeader(200)
}
