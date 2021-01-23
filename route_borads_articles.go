package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// getBoardArticles handles request with `/v1/boards/SYSOP/articles` and will return
// article list to client
func getBoardArticles(w http.ResponseWriter, r *http.Request, boardID string) {
	logger.Debugf("getBoardArticles: %v", r)
	token := getTokenFromRequest(r)
	// Check permission for board
	err := checkTokenPermission(token,
		[]permission{PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fileHeaders, err := db.ReadBoardArticleRecordsFile(boardID)

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
	}

	if _, err := w.Write(b); err != nil {
		logger.Errorf("failed to write response: %s\n", err)
	}
}

func getBoardArticlesFile(w http.ResponseWriter, r *http.Request, boardID string, filename string) {
	logger.Debugf("getBoardArticlesFile: %v", r)

	token := getTokenFromRequest(r)
	err := checkTokenPermission(token,
		[]permission{PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardID,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	buf, err := db.ReadBoardArticleFile(boardID, filename)
	if err != nil {
		logger.Errorf("read file %v/%v error: %v", boardID, filename, err)
	}

	bufStr := base64.StdEncoding.EncodeToString(buf)

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"raw": bufStr,
		},
	}

	b, err := json.MarshalIndent(responseMap, "", "  ")
	if err != nil {
		logger.Errorf("failed to marshal response map: %s\n", err)
	}

	if _, err := w.Write(b); err != nil {
		logger.Errorf("failed to write response: %s\n", err)
	}
}

func getArticleURL(boardID string, filename string) string {
	return fmt.Sprintf("https://ptt-app-dev-codingman.pichuchen.tw/bbs/%s/%s.html", boardID, filename)
}
