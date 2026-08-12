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
	mux.HandleFunc("/", delivery.notFoundHandler)
}

func (delivery *Delivery) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write(NewPathNotFoundError(r))
	if err != nil {
		delivery.logger.Errorf("handle not found handler failed: %w", err)
	}
}

func (delivery *Delivery) routeToken(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		delivery.postToken(w, r)
	default:
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write(NewMethodNotAllowedError(r, []string{http.MethodPost}))
		if err != nil {
			delivery.logger.Errorf("write NewMethodNotAllowedError error: %w", err)
		}
	}
}

func (delivery *Delivery) routeBoards(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("routeBoards: %v", r)
	switch r.Method {
	case http.MethodGet:
		delivery.getBoards(w, r)
	case http.MethodPost:
		delivery.postBoards(w, r)
	default:
		w.Header().Set("Allow", strings.Join([]string{http.MethodGet, http.MethodPost}, ","))
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write(NewMethodNotAllowedError(r, []string{http.MethodGet, http.MethodPost}))
		if err != nil {
			delivery.logger.Errorf("write NewMethodNotAllowedError error: %w", err)
		}
	}
}

func (delivery *Delivery) routePopularBoards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		delivery.getPopularBoardList(w, r)
	default:
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write(NewMethodNotAllowedError(r, []string{http.MethodGet}))
		if err != nil {
			delivery.logger.Errorf("write NewMethodNotAllowedError error: %w", err)
		}
	}
}

func (delivery *Delivery) routePopularArticles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		delivery.getPopularArticles(w, r)
	default:
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write(NewMethodNotAllowedError(r, []string{http.MethodGet}))
		if err != nil {
			delivery.logger.Errorf("write NewMethodNotAllowedError error: %w", err)
		}
	}
}

func (delivery *Delivery) routeClasses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		delivery.getClasses(w, r)
	default:
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write(NewMethodNotAllowedError(r, []string{http.MethodGet}))
		if err != nil {
			delivery.logger.Errorf("write NewMethodNotAllowedError error: %w", err)
		}
	}
}

func (delivery *Delivery) routeUsers(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("routeUsers: %v", r)
	switch r.Method {
	case http.MethodGet:
		delivery.getUsers(w, r)
	case http.MethodPost:
		delivery.postUsers(w, r)
	default:
		w.Header().Set("Allow", strings.Join([]string{http.MethodGet, http.MethodPost}, ","))
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write(NewMethodNotAllowedError(r, []string{http.MethodGet, http.MethodPost}))
		if err != nil {
			delivery.logger.Errorf("write NewMethodNotAllowedError error: %w", err)
		}
	}
}

func (delivery *Delivery) getBoards(w http.ResponseWriter, r *http.Request) {
	boardID, item, filename, err := delivery.parseBoardPath(r.URL.Path)
	if boardID == "" {
		delivery.getBoardList(w, r)
		return
	}
	if item == "information" {
		delivery.getBoardInformation(w, r, boardID)
		return
	}
	if item == "settings" {
		delivery.getBoardSettings(w, r, boardID)
		return
	}
	if item == "articles" {
		if filename == "" {
			delivery.getBoardArticles(w, r, boardID)
		} else {
			delivery.getBoardArticlesFile(w, r, boardID, filename)
		}
		return
	}
	if item == "treasures" {
		delivery.getBoardTreasures(w, r, boardID)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	_, err2 := w.Write(NewPathNotFoundError(r))
	if err2 != nil {
		delivery.logger.Errorf("write NewPathNotFoundError error: %w", err2)
	}
	delivery.logger.Noticef("board id: %v not exist but be queried, info: %v err: %v", boardID, item, err)
}

func (delivery *Delivery) postBoards(w http.ResponseWriter, r *http.Request) {
	boardID, item, filename, err := delivery.parseBoardPath(r.URL.Path)

	if boardID == "" && item == "" && filename == "" {
		delivery.createBoard(w, r)
		return
	}

	action := r.PostFormValue("action")
	if action == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write(NewPathNotFoundError(r))
		if writeErr != nil {
			delivery.logger.Errorf("postBoards write error response err: %w", writeErr)
		}
		return
	}

	if item == "articles" && boardID != "" {
		switch {
		case action == "append_comment" && filename != "":
			delivery.appendComment(w, r, boardID, filename)
			return
		case action == "forward_article" && filename != "":
			delivery.forwardArticle(w, r, boardID, filename)
			return
		case action == "add_article":
			delivery.publishPost(w, r, boardID)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	_, err2 := w.Write(NewPathNotFoundError(r))
	if err2 != nil {
		delivery.logger.Errorf("write NewPathNotFoundError error: %w", err2)
	}
	delivery.logger.Noticef("board id: %v not exist but be queried, info: %v err: %v", boardID, item, err)
}

func (delivery *Delivery) parseBoardPath(path string) (boardID string, item string, filename string, err error) {
	pathSegment := strings.Split(path, "/")
	if len(pathSegment) >= 6 {
		return pathSegment[3], pathSegment[4], pathSegment[5], nil
	}
	if len(pathSegment) == 5 {
		return pathSegment[3], pathSegment[4], "", nil
	}
	if len(pathSegment) == 4 {
		return pathSegment[3], "", "", nil
	}
	if len(pathSegment) == 3 {
		return "", "", "", nil
	}
	delivery.logger.Warningf("parseBoardPath got malform path: %v", path)
	return "", "", "", nil
}

func (delivery *Delivery) parseBoardTreasurePath(path string) (boardID string, treasuresID []string, filename string, err error) {
	pathSegment := strings.Split(path, "/")
	if len(pathSegment) == 6 {
		return pathSegment[3], []string{}, "", nil
	}
	if len(pathSegment) >= 7 {
		boardID = pathSegment[3]
		if pathSegment[len(pathSegment)-1] == "articles" {
			return boardID, pathSegment[5 : len(pathSegment)-1], "", nil
		}
		return boardID, pathSegment[5 : len(pathSegment)-2], pathSegment[len(pathSegment)-1], nil
	}
	delivery.logger.Warningf("parseBoardTreasurePath got malform path: %v", path)
	return
}

func (delivery *Delivery) parseClassPath(path string) (classID string, item string, err error) {
	pathSegment := strings.Split(path, "/")
	if len(pathSegment) == 5 {
		return pathSegment[3], pathSegment[4], nil
	}
	if len(pathSegment) == 4 {
		return pathSegment[3], "", nil
	}
	if len(pathSegment) == 3 {
		return "", "", nil
	}
	delivery.logger.Warningf("parseClassPath got malform path: %v", path)
	return "", "", nil
}

func parseUserPath(path string) (userID string, item string, itemID string, err error) {
	pathSegment := strings.Split(path, "/")
	if len(pathSegment) == 6 {
		return pathSegment[3], pathSegment[4], pathSegment[5], nil
	}
	if len(pathSegment) == 5 {
		return pathSegment[3], pathSegment[4], "", nil
	}
	if len(pathSegment) == 4 {
		return pathSegment[3], "", "", nil
	}
	return "", "", "", nil
}

func supportGuest() bool {
	return false
}
