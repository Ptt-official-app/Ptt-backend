package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (delivery *httpDelivery) buildRoute() {
	// TODO: Check IP Flowspeed
	delivery.Post("/v1/token", delivery.postToken)
	// TODO: Check IP Flowspeed
	delivery.Get("/v1/boards", delivery.getBoardList)
	delivery.Get("/v1/boards/{boardID}/information", delivery.getBoardInformation)
	delivery.Get("/v1/boards/{boardID}/articles", delivery.getBoardArticles)
	delivery.Get("/v1/boards/{boardID}/articles/{filename}", delivery.getBoardArticlesFile)
	delivery.Get("/v1/boards/{boardID}/treasures/", delivery.getBoardTreasures)
	// TODO: Check IP Flowspeed
	delivery.Get("/v1/popular-boards", delivery.getPopularBoardList)
	delivery.Get("/v1/popular-articles", delivery.getPopularArticles)
	// TODO: Check IP Flowspeed
	delivery.Get("/v1/classes/", delivery.getClasses)
	// TODO: Check IP Flowspeed
	delivery.Get("/v1/users/", delivery.getUsers)
}

func (delivery httpDelivery) Params(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func (delivery httpDelivery) Post(path string, handlerFunc http.HandlerFunc) {
	delivery.Router.HandleFunc(path, handlerFunc).Methods(http.MethodPost)
}
func (delivery httpDelivery) Get(path string, handlerFunc http.HandlerFunc) {
	delivery.Router.HandleFunc(path, handlerFunc).Methods(http.MethodGet)
}

// return a boolean value to indicate support guest account
// and using guest permission when permission insufficient
func supportGuest() bool {
	return false
}
