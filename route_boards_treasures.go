package main

import (
	// "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	// "os"
)

func getBoardTreasures(w http.ResponseWriter, r *http.Request, boardID string) {
	logger.Debugf("getBoardTreasures: %v", r)
	token := getTokenFromRequest(r)
	_, treasuresID, filename, err := parseBoardTreasurePath(r.URL.Path)

	if err != nil {
		// TODO return 400?
		logger.Warningf("parseBoardTreasurePath error: %v", err)
	}

	if filename != "" {
		// get file
		getBoardTreasuresFile(w, r, boardID, treasuresID, filename)
		return
	}

	// Check permission for board
	err = checkTokenPermission(token,
		[]permission{PermissionReadTreasureInformation},
		map[string]string{
			"board_id":    boardID,
			"treasure_id": strings.Join(treasuresID, ","),
		})
	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fileHeaders, err := db.ReadBoardTreasureRecordsFile(boardID, treasuresID)

	if err != nil {
		// The board may not contain any article
		logger.Warningf("open directory file error: %v", err)
	}

	items := []interface{}{}

	for _, f := range fileHeaders {
		m := map[string]interface{}{
			"filename": f.Filename(),
			// Bug(pichu): f.Modified time will be 0 when file is vote
			"modified_time":   f.Modified(),
			"recommend_count": f.Recommend(),
			"post_date":       f.Date(),
			"title":           f.Title(),
			"money":           fmt.Sprintf("%v", f.Money()),
			"owner":           f.Owner(),
			// "aid": ""
			"url": getArticleURL(boardID, f.Filename()),
		}
		items = append(items, m)
	}

	logger.Debugf("fh: %v", fileHeaders)

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": items,
		},
	}

	b, err := json.MarshalIndent(responseMap, "", "  ")
	if err != nil {
		logger.Errorf("failed to marshal response data: %s\n", err)
		return
	}

	if _, err := w.Write(b); err != nil {
		logger.Errorf("failed to write response: %s\n", err)
	}
}

func getBoardTreasuresFile(w http.ResponseWriter, r *http.Request, boardID string, treasuresID []string, filename string) {
	logger.Debugf("getBoardTreasuresFile %v board: %v, treasuresID: %v, filename: %v", r, boardID, treasuresID, filename)

	w.WriteHeader(http.StatusNotImplemented)
}

// parseBoardTreasurePath parse covert url path from /v1/boards/SYSOP/article to
// {SYSOP, article) or /v1/boards to {,}
func parseBoardTreasurePath(path string) (boardID string, treasuresID []string, filename string, err error) {
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
	logger.Warningf("parseBoardTreasurePath got malformed path: %v", path)

	return
}
