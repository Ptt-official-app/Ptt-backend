package http

import (
	"context"
	"encoding/json"
	"net/http"
)

func (delivery *httpDelivery) getPopularArticles(w http.ResponseWriter, r *http.Request) {
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
		w.Write(b)
		return
	}
	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": articles,
		},
	}
	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)
	return
}