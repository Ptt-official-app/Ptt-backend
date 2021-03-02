package http

import (
	"context"
	"encoding/json"
	"net/http"
)

func (delivery *httpDelivery) getPopularArticles(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getPopularArticles: %v", r)
	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": delivery.usecase.GetPopularArticles(context.Background()),
		},
	}
	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)
	return
}