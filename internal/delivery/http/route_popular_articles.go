package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

func (delivery *Delivery) getPopularArticles(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getPopularArticles: %v", r)
	articles, err := delivery.usecase.GetPopularArticles(context.Background())
	if err != nil {
		delivery.logger.Errorf("find popular articles failed: %v", err)
		m := map[string]string{
			"error":             "find_popular_articles_error",
			"error_description": "get popular articles failed",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getPopularArticles write error response err: %w", err)
		}
		return
	}
	dataList := make([]interface{}, 0, len(articles))
	for _, article := range articles {
		dataList = append(dataList, marshalArticle(article))
	}

	responseMap := map[string]interface{}{
		"data": struct {
			Items []interface{} `json:"items"`
		}{dataList},
	}
	b, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("getPopularArticles write error response err: %w", err)
	}
}

func marshalArticle(r repository.PopularArticleRecord) map[string]interface{} {
	ret := map[string]interface{}{
		"filename":        r.Filename(),
		"modified":        r.Modified(),
		"recommend_count": r.Recommend(),
		"owner":           r.Owner(),
		"date":            r.Date(),
		"title":           r.Title(),
		"money":           r.Money(),
		"board_id":        r.BoardID(),
	}

	return ret
}
