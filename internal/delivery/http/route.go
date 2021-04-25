package http

import (
	"net/http"
	"strings"
)

func (delivery *Delivery) buildRoute(mux *http.ServeMux) {
	mux.HandleFunc("/v1/token", delivery.routeToken)
	mux.HandleFunc("/v1/boards", delivery.routeBoards)
	mux.HandleFunc("/v1/boards/", delivery.routeBoards)
	mux.HandleFunc("/v1/popular-boards", delivery.routePopularBoards)
	mux.HandleFunc("/v1/popular-articles", delivery.routePopularArticles)
	mux.HandleFunc("/v1/classes/", delivery.routeClasses)
	mux.HandleFunc("/v1/users/", delivery.routeUsers)
}

func (delivery *Delivery) routeToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed
	switch r.Method {
	case http.MethodPost:
		delivery.postToken(w, r)
	}
}

// routeBoards is the handler for `/v1/boards`
func (delivery *Delivery) routeBoards(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("routeBoards: %v", r)
	// TODO: Check IP Flowspeed
	switch r.Method {
	case http.MethodGet:
		delivery.getBoards(w, r)
	case http.MethodPost:
		delivery.postBoards(w, r)
	}
}

func (delivery *Delivery) routePopularBoards(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed
	switch r.Method {
	case http.MethodGet:
		delivery.getPopularBoardList(w, r)
	}
}

// routePopularArticles a handler for `/v1/popular-articles`
func (delivery *Delivery) routePopularArticles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		delivery.getPopularArticles(w, r)
	}
}

// routeClasses is the handler for `/v1/classes`
func (delivery *Delivery) routeClasses(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed
	switch r.Method {
	case http.MethodGet:
		delivery.getClasses(w, r)
	}
}

// routeClasses is the handler for `/v1/users`
func (delivery *Delivery) routeUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed
	delivery.logger.Debugf("routeUsers: %v", r)
	switch r.Method {
	case http.MethodGet:
		delivery.getUsers(w, r)
	case http.MethodPost:
		delivery.postUsers(w, r)
	}
}

// getBoards is the handler for `/v1/boards` with GET method
func (delivery *Delivery) getBoards(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getBoards: %v", r)
	boardID, item, filename, err := delivery.parseBoardPath(r.URL.Path)
	if boardID == "" {
		delivery.getBoardList(w, r)
		return
	}
	// get single board
	if item == "information" {
		delivery.getBoardInformation(w, r, boardID)
		return
	} else if item == "settings" {
		delivery.getBoardSettings(w, r, boardID)
		return
	} else if item == "articles" {
		if filename == "" {
			delivery.getBoardArticles(w, r, boardID)
		} else {
			delivery.getBoardArticlesFile(w, r, boardID, filename)
		}
		return
	} else if item == "treasures" {
		delivery.getBoardTreasures(w, r, boardID)
		return
	}

	// 404
	w.WriteHeader(http.StatusNotFound)

	delivery.logger.Noticef("board id: %v not exist but be queried, info: %v err: %v", boardID, item, err)
}

// postBoards is the handler for `/v1/boards` with POST method
func (delivery *Delivery) postBoards(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("postBoards: %v", r)
	boardID, item, filename, err := delivery.parseBoardPath(r.URL.Path)

	action := r.PostFormValue("action")
	if action == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if item == "articles" && boardID != "" {
		if action == "append_comment" && filename != "" {
			delivery.appendComment(w, r, boardID, filename)
			return
		} else if action == "forward_article" && filename != "" {
			delivery.forwardArticle(w, r, boardID, filename)
		} else if action == "add_article" {
			delivery.publishPost(w, r, boardID)
			return
		}
	}

	// 404
	w.WriteHeader(http.StatusNotFound)

	delivery.logger.Noticef("board id: %v not exist but be queried, info: %v err: %v", boardID, item, err)
}

// parseBoardPath covert url path from /v1/boards/SYSOP/article to
// {SYSOP, article) or /v1/boards to {,}
func (delivery *Delivery) parseBoardPath(path string) (boardID string, item string, filename string, err error) {
	pathSegment := strings.Split(path, "/")

	if len(pathSegment) >= 6 {
		// /{{version}}/boards/{{class_id}}/{{item}}/{{filename}}
		boardID = pathSegment[3]
		item = pathSegment[4]
		filename = pathSegment[5]
		return
	} else if len(pathSegment) == 5 {
		// /{{version}}/boards/{{class_id}}/{{item}}
		boardID = pathSegment[3]
		item = pathSegment[4]
		return
	} else if len(pathSegment) == 4 {
		// /{{version}}/boards/{{class_id}}
		boardID = pathSegment[3]
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
func (delivery *Delivery) parseBoardTreasurePath(path string) (boardID string, treasuresID []string, filename string, err error) {
	pathSegment := strings.Split(path, "/")

	if len(pathSegment) == 6 {
		// /{{version}}/boards/{{board_id}}/treasures/articles
		boardID = pathSegment[3]
		treasuresID = []string{}
		filename = ""
		return
	} else if len(pathSegment) >= 7 {
		// /{{version}}/boards/{{board_id}}/treasures/{{treasures_id ... }}/articles
		// or
		// /{{version}}/boards/{{board_id}}/treasures/{{treasures_id ... }}/articles/{{filename}}
		boardID = pathSegment[3]
		if pathSegment[len(pathSegment)-1] == "articles" {
			treasuresID = pathSegment[5 : len(pathSegment)-1]
			filename = ""
		} else {
			treasuresID = pathSegment[5 : len(pathSegment)-2]
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
func (delivery *Delivery) parseClassPath(path string) (classID string, item string, err error) {
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

func parseUserPath(path string) (userID string, item string, itemID string, err error) {
	pathSegment := strings.Split(path, "/")
	if len(pathSegment) == 6 {
		// /{{version}}/users/{{user_id}}/{{item}}/{{itemID}}
		return pathSegment[3], pathSegment[4], pathSegment[5], nil
	} else if len(pathSegment) == 5 {
		// /{{version}}/users/{{user_id}}/{{item}}
		return pathSegment[3], pathSegment[4], "", nil
	} else if len(pathSegment) == 4 {
		// /{{version}}/users/{{user_id}}
		return pathSegment[3], "", "", nil
	}
	return "", "", "", nil
}

// return a boolean value to indicate support guest account
// and using guest permission when permission insufficient
func supportGuest() bool {
	return false
}
