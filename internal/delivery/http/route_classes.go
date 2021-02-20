package http

import (
	"context"
	// "github.com/PichuChen/go-bbs"
	// "github.com/PichuChen/go-bbs/crypt"
	// "log"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (delivery *httpDelivery) getClass(w http.ResponseWriter, r *http.Request) {

	seg := strings.Split(r.URL.Path, "/")

	classID := "0"
	if len(seg) > 2 {
		classID = seg[3]
	}
	delivery.logger.Informationalf("user get class: %v", classID)

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

	_, _ = w.Write(b)

}

// getClasses HandleFunc handles path start with `/v1/classes`
// and pass requests to next handle function
func (delivery *httpDelivery) getClasses(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getClasses: %v", r)
	classID, item, err := delivery.parseClassPath(r.URL.Path)
	delivery.logger.Noticef("query class: %v item: %v err: %v", classID, item, err)
	if classID == "" {
		getClassesWithoutClassID(w, r)
		return
	}
	delivery.getClassesList(w, r, classID)

	// // get single board
	// if item == "information" {
	// 	getBoardInformation(w, r, boardId)
	// 	return
	// }

}

// getClassesWithoutClassID handles path don't contain item after class id
// eg: `/v1/classes`, it will redirect Client to `/v1/classes/1` which is
// root class by default.
func getClassesWithoutClassID(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/v1/classes/1", 301)
}

// getClassesList handle path with class id and will return boards and classes
// under this class.
// TODO: What should we return when target class not found?
func (delivery *httpDelivery) getClassesList(w http.ResponseWriter, r *http.Request, classID string) {
	delivery.logger.Debugf("getClassesList: %v", r)

	token := delivery.getTokenFromRequest(r)
	userID, err := delivery.usecase.GetUserIDFromToken(token)
	if err != nil {
		// user permission error
		// Support Guest?
		userID = "guest" // TODO: use const variable

		if !supportGuest() {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"token_invalid"}`))
			return
		}
	}

	boards := delivery.usecase.GetClasses(context.Background(), userID, classID)

	dataList := []interface{}{}
	for bid, b := range boards {
		m := marshalBoardHeader(b)
		if b.IsClass() {
			m["id"] = fmt.Sprintf("%v", bid+1)
		}
		dataList = append(dataList, m)
	}

	responseMap := map[string]interface{}{
		"data": dataList,
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)
}
