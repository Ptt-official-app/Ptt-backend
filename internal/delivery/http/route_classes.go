package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// getClasses HandleFunc handles path start with `/v1/classes`
// and pass requests to next handle function
func (delivery *Delivery) getClasses(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getClasses: %v", r)
	classID, item, err := delivery.parseClassPath(r.URL.Path)
	delivery.logger.Noticef("query class: %v item: %v err: %v", classID, item, err)
	if classID == "" {
		getClassesWithoutClassID(w, r)
		return
	}
	delivery.getClassesList(w, r, classID)
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
func (delivery *Delivery) getClassesList(w http.ResponseWriter, r *http.Request, classID string) {
	delivery.logger.Debugf("getClassesList: %v", r)

	token := delivery.getTokenFromRequest(r)
	userID, err := delivery.usecase.GetUserIDFromToken(token)
	if err != nil {
		userID = "guest" // TODO: use const variable
		// user permission error
		// Support Guest?
		if !supportGuest() {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = w.Write([]byte(`{"error":"token_invalid"}`))
			if err != nil {
				delivery.logger.Errorf("getClassesList write error response err: %w", err)
			}
			return
		}
	}

	if classID != "1" {
		i, err := strconv.Atoi(classID)
		if err == nil {
			classID = fmt.Sprintf("%v", i-ClassIDBase)
		}
	}

	boards, err := delivery.usecase.GetClasses(context.Background(), userID, classID)
	if err != nil {
		delivery.logger.Warningf("find classes %s failed: %v", classID, err)
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_classes_error",
			"error_description": "get classes for " + classID + " failed",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getClassesList write error response err: %w", err)
		}
		return
	}

	dataList := []interface{}{}
	for bid, b := range boards {
		m := marshalBoardHeader(b)
		if b.IsClass() {
			m["id"] = fmt.Sprintf("%v", bid+1+ClassIDBase)
		}
		dataList = append(dataList, m)
	}

	responseMap := map[string]interface{}{
		"data": dataList,
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("getClassesList write success response err: %w", err)
	}
}
