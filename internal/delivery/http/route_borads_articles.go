package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PichuChen/go-bbs"
)

// getBoardArticles handles request with `/v1/boards/SYSOP/articles` and will return
// article list to client
func (delivery *httpDelivery) getBoardArticles(w http.ResponseWriter, r *http.Request, boardId string) {
	delivery.logger.Debugf("getBoardArticles: %v", r)
	token := delivery.getTokenFromRequest(r)
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
	fileHeaders, err = delivery.boardRepo.GetBoardArticleRecords(context.Background(), boardId)
	if err != nil {
		delivery.logger.Warningf("open directory file error: %v", err)
		// The board may not contain any article
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
			"url": getArticleURL(boardId, f.Filename()),
		}
		items = append(items, m)
	}
	delivery.logger.Debugf("fh: %v", fileHeaders)

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": items,
		},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)

}

func (delivery *httpDelivery) getBoardArticlesFile(w http.ResponseWriter, r *http.Request, boardId string, filename string) {
	delivery.logger.Debugf("getBoardArticlesFile: %v", r)

	token := delivery.getTokenFromRequest(r)
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

	buf, err := delivery.boardRepo.GetBoardArticle(context.Background(), boardId, filename)
	if err != nil {
		delivery.logger.Errorf("read file %v/%v error: %v", boardId, filename, err)
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
