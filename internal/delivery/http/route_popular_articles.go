package http

import (
	"context"
	"encoding/json"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
	"net/http"
)

func (delivery *httpDelivery) getPopularArticles(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getPopularArticles: %v", r)
	// TODO: replace GetBoardArticles method when GetPopularArticles usecase is ready
	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": delivery.usecase.GetBoardArticles(context.Background(), "ALLPOST", &usecase.ArticleSearchCond{}),
		},
	}
	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)
	return
}