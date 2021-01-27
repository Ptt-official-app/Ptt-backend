package main

import "github.com/julienschmidt/httprouter"

func buildRoute(r *httprouter.Router) {

	r.POST("/v1/token", postToken)
	r.GET("/v1/boards", getBoards)
	r.GET("/v1/classes/", getClasses)
	r.GET("/v1/users/", getUsers)

}
