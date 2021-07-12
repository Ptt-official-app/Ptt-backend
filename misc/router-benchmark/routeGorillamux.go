package router_benchmark

import (
	"github.com/gorilla/mux"
	"net/http"
)

// getBoards is the handler for `/v1/boards` with GET method
func getBoardArticles_gorillamux(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	boardID := params["boardID"]
	_ = boardID
	w.WriteHeader(200)
}
