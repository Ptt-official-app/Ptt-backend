package rest

import (
	"encoding/json"
	"net/http"
	"strings"
)

func buildRoute(r *http.ServeMux) {

	r.HandleFunc("/v1/token", routeToken)
	r.HandleFunc("/v1/boards", routeBoards)
	r.HandleFunc("/v1/boards/", routeBoards)
	r.HandleFunc("/v1/classes/", routeClasses)
	r.HandleFunc("/v1/users/", routeUsers)
}

func routeClass(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getClass(w, r)
		return
	}

}

func getClass(w http.ResponseWriter, r *http.Request) {

	seg := strings.Split(r.URL.Path, "/")

	classId := "0"
	if len(seg) > 2 {
		classId = seg[3]
	}
	logger.Informationalf("user get class: %v", classId)

	list := []interface{}{}

	c := map[string]interface{}{
		"id":             1,
		"type":           "class",
		"title":          "title",
		"number_of_user": 3,
		"moderators": []string{
			"SYSOP",
			"pichu",
		},
	}
	list = append(list, c)

	m := map[string]interface{}{
		"data": list,
	}
	b, _ := json.MarshalIndent(m, "", "  ")

	w.Write(b)

}

// return a boolean value to indicate support guest account
// and using guest permission when permission insufficient
func supportGuest() bool {
	return false
}
