package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

// getBoardArticles handles request with `/v1/boards/SYSOP/articles` and will return
// article list to client
func (delivery *httpDelivery) getBoardArticles(w http.ResponseWriter, r *http.Request, boardId string) {
	delivery.logger.Debugf("getBoardArticles: %v", r)
	token := delivery.getTokenFromRequest(r)
	// Check permission for board
	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var recommendCountGreaterEqual, recommendCountLessEqual int
	var recommendCountGreaterEqualIsSet, recommendCountLessEqualIsSet bool
	queryParam := r.URL.Query()
	recommendCountGreaterEqualParam := queryParam.Get("recommend_count_ge")
	recommendCountGreaterEqualIsSet = recommendCountGreaterEqualParam != ""
	recommendCountGreaterEqual, err = strconv.Atoi(recommendCountGreaterEqualParam)

	if err != nil && recommendCountGreaterEqualIsSet {
		delivery.logger.Errorf("recommend_count_ge should be integer")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recommendCountLessEqualParam := queryParam.Get("recommend_count_le")
	recommendCountLessEqualIsSet = recommendCountLessEqualParam != ""
	recommendCountLessEqual, err = strconv.Atoi(recommendCountLessEqualParam)
	if err != nil && recommendCountLessEqualIsSet {
		delivery.logger.Errorf("recommend_count_le should be integer")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	searchCond := &usecase.ArticleSearchCond{
		Title:                           queryParam.Get("title_contain"),
		Author:                          queryParam.Get("author"),
		RecommendCountGreaterEqual:      recommendCountGreaterEqual,
		RecommendCountLessEqual:         recommendCountLessEqual,
		RecommendCountGreaterEqualIsSet: recommendCountGreaterEqualIsSet,
		RecommendCountLessEqualIsSet:    recommendCountLessEqualIsSet,
	}

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": delivery.usecase.GetBoardArticles(context.Background(), boardId, searchCond),
		},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)

}

func (delivery *httpDelivery) getBoardArticlesFile(w http.ResponseWriter, r *http.Request, boardId string, filename string) {
	delivery.logger.Debugf("getBoardArticlesFile: %v", r)

	token := delivery.getTokenFromRequest(r)
	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	buf, err := delivery.usecase.GetBoardArticle(context.Background(), boardId, filename)
	if err != nil {
		delivery.logger.Errorf("failed to get board article: %s", err)
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
