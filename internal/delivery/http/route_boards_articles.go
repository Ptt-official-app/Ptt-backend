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
func (delivery *Delivery) getBoardArticles(w http.ResponseWriter, r *http.Request, boardID string) {
	delivery.logger.Debugf("getBoardArticles: %v", r)
	token := delivery.getTokenFromRequest(r)
	// Check permission for board
	err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionReadBoardInformation}, map[string]string{
		"board_id": boardID,
	})
	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write(NewNoPermissionForReadBoardArticlesError(r, boardID))
		if err != nil {
			delivery.logger.Errorf("write NewNoPermissionForReadBoardArticlesError error: %w", err)
		}
		return
	}

	var recommendCountGreater, recommendCountLess int
	var recommendCountGreaterEqualIsSet, recommendCountLessEqualIsSet bool
	queryParam := r.URL.Query()
	getRecommendCount := func(name string) (*int, error) {
		recommendCountParam := queryParam.Get(name)
		if recommendCountParam == "" {
			return nil, nil
		}
		recommendCount, err := strconv.Atoi(recommendCountParam)
		if err != nil {
			return nil, err
		}
		return &recommendCount, nil
	}

	recommendCountGreaterEqual, err := getRecommendCount("recommend_count_ge")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(NewParameterShouldBeIntegerError(r, "recommend_count_ge"))
		if err != nil {
			delivery.logger.Errorf("got error %w when write ParameterShouldBeIntegerError", err)
		}
		return
	}

	recommendCountGreaterThan, err := getRecommendCount("recommend_count_gt")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(NewParameterShouldBeIntegerError(r, "recommend_count_gt"))
		if err != nil {
			delivery.logger.Errorf("got error %w when write ParameterShouldBeIntegerError", err)
		}
		return
	}

	recommendCountLessEqual, err := getRecommendCount("recommend_count_le")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(NewParameterShouldBeIntegerError(r, "recommend_count_le"))
		if err != nil {
			delivery.logger.Errorf("got error %w when write ParameterShouldBeIntegerError", err)
		}
		return
	}

	recommendCountLessThan, err := getRecommendCount("recommend_count_lt")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(NewParameterShouldBeIntegerError(r, "recommend_count_lt"))
		if err != nil {
			delivery.logger.Errorf("got error %w when write ParameterShouldBeIntegerError", err)
		}
		return
	}

	if recommendCountLessThan == nil && recommendCountLessEqual == nil {
		recommendCountLessEqualIsSet = false
	} else if recommendCountLessThan != nil && recommendCountLessEqual != nil {
		recommendCountLessEqualIsSet = true
		recommendCountLess = *recommendCountLessThan - 1
		if recommendCountLess > *recommendCountLessEqual {
			recommendCountLess = *recommendCountLessEqual
		}
	} else if recommendCountLessThan != nil {
		recommendCountLessEqualIsSet = true
		recommendCountLess = *recommendCountLessThan - 1
	} else {
		recommendCountLessEqualIsSet = true
		recommendCountLess = *recommendCountLessEqual
	}

	if recommendCountGreaterThan == nil && recommendCountGreaterEqual == nil {
		recommendCountGreaterEqualIsSet = false
	} else if recommendCountGreaterThan != nil && recommendCountGreaterEqual != nil {
		recommendCountGreaterEqualIsSet = true
		recommendCountGreater = *recommendCountGreaterThan + 1
		if recommendCountGreater < *recommendCountGreaterEqual {
			recommendCountGreater = *recommendCountGreaterEqual
		}
	} else if recommendCountGreaterThan != nil {
		recommendCountGreaterEqualIsSet = true
		recommendCountGreater = *recommendCountGreaterThan + 1
	} else {
		recommendCountGreaterEqualIsSet = true
		recommendCountGreater = *recommendCountGreaterEqual
	}

	searchCond := &usecase.ArticleSearchCond{
		Title:                           queryParam.Get("title_contain"),
		Author:                          queryParam.Get("author"),
		RecommendCountGreaterEqual:      recommendCountGreater,
		RecommendCountLessEqual:         recommendCountLess,
		RecommendCountGreaterEqualIsSet: recommendCountGreaterEqualIsSet,
		RecommendCountLessEqualIsSet:    recommendCountLessEqualIsSet,
	}

	items := []interface{}{}
	articles := delivery.usecase.GetBoardArticles(context.Background(), boardID, searchCond)

	// Articles to output format, please refer: https://docs.google.com/document/d/18DsZOyrlr5BIl2kKxZH7P2QxFLG02xL2SO0PzVHVY3k/edit#heading=h.bnhpxsiwnbey
	for _, a := range articles {
		item := map[string]interface{}{
			"filename":        a.Filename(),
			"modified_time":   a.Modified(),
			"recommend_count": a.Recommend(),
			"post_date":       a.Date(),
			"title":           a.Title(),
			"money":           a.Money(),
			"owner":           a.Owner(),
			// TODO: generate aid and url
			// "aid": a.Aid(),
		}
		u, ok := delivery.usecase.(usecase.SupportWebUsecase)
		if ok {
			item["url"] = u.GetArticleURL(boardID, a.Filename())
		}

		items = append(items, item)
	}

	responseMap := map[string]interface{}{
		"data": map[string]interface{}{
			"items": items,
		},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("getBoardArticles write success response err: %w", err)
	}
}

func (delivery *Delivery) getBoardArticlesFile(w http.ResponseWriter, r *http.Request, boardID string, filename string) {
	delivery.logger.Debugf("getBoardArticlesFile: %v", r)
	ctx := context.Background()

	token := delivery.getTokenFromRequest(r)
	err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionReadBoardInformation}, map[string]string{
		"board_id": boardID,
	})
	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	buf, err := delivery.usecase.GetBoardArticle(ctx, boardID, filename)
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
	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("getBoardArticlesFile write success response err: %w", err)
	}
}
