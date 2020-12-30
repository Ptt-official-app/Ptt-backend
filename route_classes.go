package main

import (
	// "github.com/PichuChen/go-bbs"
	// "github.com/PichuChen/go-bbs/crypt"
	// "log"
	"net/http"
	// "strings"
)

func routeClasses(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed

	if r.Method == "GET" {
		getClasses(w, r)
		return
	}

}

func getClasses(w http.ResponseWriter, r *http.Request) {
	logger.Criticalf("get classes not implement")

	w.WriteHeader(http.StatusNotImplemented)

}
