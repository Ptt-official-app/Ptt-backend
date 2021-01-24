package main

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

// routeClasses is the handler for `/v1/classes`
func routeClasses(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed

	if r.Method == "GET" {
		getClasses(w, r)
		return
	}

}

// getClasses HandleFunc handles path start with `/v1/classes`
// and pass requests to next handle function
func getClasses(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("getClasses: %v", r)
	classId, item, err := parseClassPath(r.URL.Path)
	logger.Noticef("query class: %v item: %v err: %v", classId, item, err)
	if classId == "" {
		getClassesWithoutClassId(w, r)
		return
	}
	getClassesList(w, r, classId)
	return

	// // get single board
	// if item == "information" {
	// 	getBoardInformation(w, r, boardId)
	// 	return
	// }

}

// getClassesWithoutClassId handles path don't contain item after class id
// eg: `/v1/classes`, it will redirect Client to `/v1/classes/1` which is
// root class by default.
func getClassesWithoutClassId(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/v1/classes/1", 301)
}

// getClassesList handle path with class id and will return boards and classes
// under this class.
// TODO: What should we return when target class not found?
func getClassesList(w http.ResponseWriter, r *http.Request, classId string) {
	logger.Debugf("getClassesList: %v", r)

	token := getTokenFromRequest(r)
	userId, err := getUserIdFromToken(token)
	if err != nil {
		// user permission error
		// Support Guest?
		if !supportGuest() {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"token_invalid"}`))
			return
		} else {
			userId = "guest" // TODO: use const variable
		}
	}

	dataList := []interface{}{}
	for bid, b := range boardRepo.GetBoards(context.Background()) {
		// TODO: Show Board by user level
		if !shouldShowOnUserLevel(b, userId) {
			continue
		}
		if b.ClassId() != classId {
			continue
		}
		jb, _ := json.Marshal(b)
		logger.Debugf("marshal class board: %v", string(jb))
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

// parseClassPath covert url path from /v1/classes/1/information to
// {1, information) or /v1/classes to {,}
func parseClassPath(path string) (classId string, item string, err error) {
	pathSegment := strings.Split(path, "/")
	if len(pathSegment) == 5 {
		// /{{version}}/classes/{{class_id}}/{{item}}
		return pathSegment[3], pathSegment[4], nil
	} else if len(pathSegment) == 4 {
		// /{{version}}/classes/{{class_id}}
		return pathSegment[3], "", nil
	} else if len(pathSegment) == 3 {
		// /{{version}}/classes
		return "", "", nil
	}
	logger.Warningf("parseClassPath got malform path: %v", path)
	return "", "", nil

}
