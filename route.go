package main

import (
	"net/http"
)

func buildRoute(r *http.ServeMux) {

	r.HandleFunc("/v1/token", routeToken)
	r.HandleFunc("/v1/boards", routeBoards)
	r.HandleFunc("/v1/boards/", routeBoards)
	r.HandleFunc("/v1/classes/", routeClasses)
	r.HandleFunc("/v1/users/", routeUsers)
}
