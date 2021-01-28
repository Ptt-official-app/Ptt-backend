package http

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (delivery *httpDelivery) routeClass(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		delivery.getClass(w, r)
		return
	}

}

func (delivery *httpDelivery) getClass(w http.ResponseWriter, r *http.Request) {

	seg := strings.Split(r.URL.Path, "/")

	classId := "0"
	if len(seg) > 2 {
		classId = seg[3]
	}
	delivery.logger.Informationalf("user get class: %v", classId)

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
