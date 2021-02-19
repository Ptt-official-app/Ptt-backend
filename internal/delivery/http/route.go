package http

import (
	"net/http"
	"strings"
)

func (delivery *httpDelivery) buildRoute(mux *http.ServeMux) {
	mux.HandleFunc("/v1/token", delivery.routeToken)
	mux.HandleFunc("/v1/boards", delivery.routeBoards)
	mux.HandleFunc("/v1/boards/", delivery.routeBoards)
	mux.HandleFunc("/v1/popular-boards", delivery.routePopularBoards)
	mux.HandleFunc("/v1/classes/", delivery.routeClasses)
	mux.HandleFunc("/v1/users/", delivery.routeUsers)
}

func (delivery *httpDelivery) routeToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed
	switch r.Method {
	case http.MethodPost:
		delivery.postToken(w, r)
	}
}

func (delivery *httpDelivery) routeClass(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		delivery.getClass(w, r)
	}
}

// routeBoards is the handler for `/v1/boards`
func (delivery *httpDelivery) routeBoards(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("routeBoards: %v", r)
	// TODO: Check IP Flowspeed
	switch r.Method {
	case http.MethodGet:
		delivery.getBoards(w, r)
	}
}


func (delivery *httpDelivery) routePopularBoards(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed
	switch r.Method {
	case http.MethodGet:
		delivery.getPopularBoardList(w, r)
	}
}

// routeClasses is the handler for `/v1/classes`
func (delivery *httpDelivery) routeClasses(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed
	switch r.Method {
	case http.MethodGet:
		delivery.getClasses(w, r)
	}
}

func (delivery *httpDelivery) routeUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed
	switch r.Method {
	case http.MethodGet:
		delivery.getUsers(w, r)
	}
}

// getBoards is the handler for `/v1/boards` with GET method
func (delivery *httpDelivery) getBoards(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getBoards: %v", r)
	boardId, item, filename, err := delivery.parseBoardPath(r.URL.Path)
	if boardId == "" {
		delivery.getBoardList(w, r)
		return
	}
	// get single board
	if item == "information" {
		delivery.getBoardInformation(w, r, boardId)
		return
	} else if item == "articles" {
		if filename == "" {
			delivery.getBoardArticles(w, r, boardId)
		} else {
			delivery.getBoardArticlesFile(w, r, boardId, filename)
		}
		return
	} else if item == "treasures" {
		delivery.getBoardTreasures(w, r, boardId)
		return
	}

	// 404
	w.WriteHeader(http.StatusNotFound)

	delivery.logger.Noticef("board id: %v not exist but be queried, info: %v err: %v", boardId, item, err)
}

// parseBoardPath covert url path from /v1/boards/SYSOP/article to
// {SYSOP, article) or /v1/boards to {,}
func (delivery *httpDelivery) parseBoardPath(path string) (boardId string, item string, filename string, err error) {
	pathSegment := strings.Split(path, "/")

	if len(pathSegment) >= 6 {
		// /{{version}}/boards/{{class_id}}/{{item}}/{{filename}}
		boardId = pathSegment[3]
		item = pathSegment[4]
		filename = pathSegment[5]
		return
	} else if len(pathSegment) == 5 {
		// /{{version}}/boards/{{class_id}}/{{item}}
		boardId = pathSegment[3]
		item = pathSegment[4]
		return
	} else if len(pathSegment) == 4 {
		// /{{version}}/boards/{{class_id}}
		boardId = pathSegment[3]
		return
	} else if len(pathSegment) == 3 {
		// /{{version}}/boards
		// Should not be reach...
		return
	}
	delivery.logger.Warningf("parseBoardPath got malform path: %v", path)
	return
}

// parseBoardTreasurePath parse covert url path from /v1/boards/SYSOP/article to
// {SYSOP, article) or /v1/boards to {,}
func (delivery *httpDelivery) parseBoardTreasurePath(path string) (boardId string, treasuresId []string, filename string, err error) {
	pathSegment := strings.Split(path, "/")

	if len(pathSegment) == 6 {
		// /{{version}}/boards/{{board_id}}/treasures/articles
		boardId = pathSegment[3]
		treasuresId = []string{}
		filename = ""
		return
	} else if len(pathSegment) >= 7 {
		// /{{version}}/boards/{{board_id}}/treasures/{{treasures_id ... }}/articles
		// or
		// /{{version}}/boards/{{board_id}}/treasures/{{treasures_id ... }}/articles/{{filename}}
		boardId = pathSegment[3]
		if pathSegment[len(pathSegment)-1] == "articles" {
			treasuresId = pathSegment[5 : len(pathSegment)-1]
			filename = ""
		} else {
			treasuresId = pathSegment[5 : len(pathSegment)-2]
			filename = pathSegment[len(pathSegment)-1]
		}
		return
	}
	// should not be reached
	delivery.logger.Warningf("parseBoardTreasurePath got malform path: %v", path)
	return
}

// parseClassPath covert url path from /v1/classes/1/information to
// {1, information) or /v1/classes to {,}
func (delivery *httpDelivery) parseClassPath(path string) (classId string, item string, err error) {
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
	delivery.logger.Warningf("parseClassPath got malform path: %v", path)
	return "", "", nil
}

func parseUserPath(path string) (userId string, item string, err error) {
	pathSegment := strings.Split(path, "/")
	// /{{version}}/users/{{user_id}}/{{item}}
	if len(pathSegment) == 4 {
		// /{{version}}/users/{{user_id}}
		return pathSegment[3], "", nil
	}
	return pathSegment[3], pathSegment[4], nil
}

// return a boolean value to indicate support guest account
// and using guest permission when permission insufficient
func supportGuest() bool {
	return false
}
