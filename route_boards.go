package main

import (
	"encoding/json"
	// "fmt"
	// "github.com/PichuChen/go-bbs"
	// "github.com/PichuChen/go-bbs/crypt"
	// "github.com/dgrijalva/jwt-go"
	"net/http"
	// "strings"
)

func routeBoards(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed

	if r.Method == "GET" {
		getBoards(w, r)
		return
	}

}

func getBoards(w http.ResponseWriter, r *http.Request) {

	// TODO: Check JWT

	// TODO: Get user Level

	// TODO: Show Board by user level

	dataList := []interface{}{}
	for _, b := range boardHeader {
		dataList = append(dataList, b)
	}

	responseMap := map[string]interface{}{
		"data": dataList,
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)

}
