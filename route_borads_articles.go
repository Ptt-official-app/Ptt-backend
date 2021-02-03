package main

import (
	"github.com/PichuChen/go-bbs"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// getBoardArticles handles request with `/v1/boards/SYSOP/articles` and will return
// article list to client
func getBoardArticles(w http.ResponseWriter, r *http.Request, boardId string) {
	logger.Debugf("getBoardArticles: %v", r)
	token := getTokenFromRequest(r)
	// Check permission for board
	err := checkTokenPermission(token,
		[]permission{PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var fileHeaders []bbs.ArticleRecord
	fileHeaders, err = db.ReadBoardArticleRecordsFile(boardId)
	if err != nil {
		logger.Warningf("open directory file error: %v", err)
		// The board may not contain any article
	}

	var items []interface{}
	var articles []bbs.ArticleRecord
	var recommendCountGe int
	queryParam := r.URL.Query()
	title := queryParam.Get("title_contain")
	author := queryParam.Get("author")
	recommendCountParam := queryParam.Get("recommend_count_ge")
	recommendCountGe, err = strconv.Atoi(recommendCountParam)
	if err != nil && recommendCountParam != "" {
		logger.Errorf("recommend count should be integer")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(title)) > 0 ||
		len(strings.TrimSpace(author)) > 0 ||
		recommendCountGe != 0 {
		articles = searchArticles(fileHeaders, title, author, recommendCountGe)
	} else {
		articles = fileHeaders
	}

	for _, f := range articles {
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
			"url": getArticleURL(boardId, f.Filename()),
		}
		items = append(items, m)
	}
	logger.Debugf("fh: %v", fileHeaders)

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": items,
		},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)

}

func searchArticles(fileHeaders []bbs.ArticleRecord, title, author string, recommendCountGe int) []bbs.ArticleRecord {
	var targetArticles []bbs.ArticleRecord

	for _, f := range fileHeaders {
		if strings.Contains(strings.ToLower(f.Title()), strings.ToLower(title)) &&
			strings.Contains(strings.ToLower(f.Owner()), strings.ToLower(author)) &&
			((recommendCountGe == 0) ||
				(recommendCountGe > 0 && f.Recommend() >= recommendCountGe) ||
				(recommendCountGe < 0 && f.Recommend() <= recommendCountGe)) {
			targetArticles = append(targetArticles, f)
		}
	}
	return targetArticles
}

func getBoardArticlesFile(w http.ResponseWriter, r *http.Request, boardId string, filename string) {
	logger.Debugf("getBoardArticlesFile: %v", r)

	token := getTokenFromRequest(r)
	err := checkTokenPermission(token,
		[]permission{PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	buf, err := db.ReadBoardArticleFile(boardId, filename)
	if err != nil {
		logger.Errorf("read file %v/%v error: %v", boardId, filename, err)
	}

	bufStr := base64.StdEncoding.EncodeToString(buf)

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"raw": bufStr,
		},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)
}

func getArticleURL(boardId string, filename string) string {
	return fmt.Sprintf("https://ptt-app-dev-codingman.pichuchen.tw/bbs/%s/%s.html", boardId, filename)
}
