package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func NewURLValuesParser(val url.Values, gek, gtk, lek, ltk string) URLValuesParser {
	return URLValuesParser{
		Values:          val,
		greaterEqualKey: gek,
		greaterThanKey:  gtk,
		lessEqualKey:    lek,
		lessThanKey:     ltk,
		errorCalled:     false,
	}
}

type URLValuesParser struct {
	url.Values
	greaterEqualKey string
	greaterThanKey  string
	lessEqualKey    string
	lessThanKey     string

	greaterEqual *int
	lessEqual    *int
	errorCalled  bool
}

func (parser *URLValuesParser) getRecommendCount(name string) (*int, error) {
	recommendCountParam := parser.Values.Get(name)
	if recommendCountParam == "" {
		return nil, nil
	}
	recommendCount, err := strconv.Atoi(recommendCountParam)
	if err != nil {
		return nil, err
	}
	return &recommendCount, nil
}

func (parser *URLValuesParser) Error() (string, error) {
	parser.errorCalled = true
	greaterEqual, err := parser.getRecommendCount(parser.greaterEqualKey)
	if err != nil {
		return parser.greaterEqualKey, err
	}

	greaterThan, err := parser.getRecommendCount(parser.greaterThanKey)
	if err != nil {
		return parser.greaterThanKey, err
	}

	// check the intersection of greater
	if greaterEqual == nil && greaterThan == nil {
	} else if greaterEqual != nil && greaterThan != nil {
		// x > 10 === x >= 11
		*greaterThan++
		parser.greaterEqual = greaterThan
		if *parser.greaterEqual < *greaterEqual {
			parser.greaterEqual = greaterEqual
		}
	} else if greaterThan != nil {
		*greaterThan++
		parser.greaterEqual = greaterThan
	} else {
		parser.greaterEqual = greaterEqual
	}

	lessEqual, err := parser.getRecommendCount(parser.lessEqualKey)
	if err != nil {
		return parser.lessEqualKey, err
	}

	lessThan, err := parser.getRecommendCount(parser.lessThanKey)
	if err != nil {
		return parser.lessThanKey, err
	}

	// check the intersection of less
	if lessEqual == nil && lessThan == nil {
	} else if lessEqual != nil && lessThan != nil {
		// x < 10 === x <= 9
		*lessThan--
		parser.lessEqual = lessThan
		if *parser.lessEqual > *lessEqual {
			parser.lessEqual = lessEqual
		}
	} else if lessThan != nil {
		*lessThan--
		parser.lessEqual = lessThan
	} else {
		parser.lessEqual = lessEqual
	}
	return "", nil
}

func (parser *URLValuesParser) GetGreaterEqual() (int, bool) {
	if !parser.errorCalled {
		panic("URLValuesParser.Error function must be called first")
	}
	if parser.greaterEqual == nil {
		return 0, false
	}
	return *parser.greaterEqual, true
}

func (parser *URLValuesParser) GetLessEqual() (int, bool) {
	if !parser.errorCalled {
		panic("URLValuesParser.Error function must be called first")
	}
	if parser.lessEqual == nil {
		return 0, false
	}
	return *parser.lessEqual, true
}

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

	queryParam := r.URL.Query()
	parser := NewURLValuesParser(queryParam, "recommend_count_ge", "recommend_count_gt", "recommend_count_le", "recommend_count_lt")
	errorKey, err := parser.Error()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(NewParameterShouldBeIntegerError(r, errorKey))
		if err != nil {
			delivery.logger.Errorf("got error %w when write ParameterShouldBeIntegerError", err)
		}
		return
	}

	recommendCountGreaterEqual, recommendCountGreaterEqualIsSet := parser.GetGreaterEqual()
	recommendCountLessEqual, recommendCountLessEqualIsSet := parser.GetLessEqual()

	searchCond := &usecase.ArticleSearchCond{
		Title:                           queryParam.Get("title_contain"),
		Author:                          queryParam.Get("author"),
		RecommendCountGreaterEqual:      recommendCountGreaterEqual,
		RecommendCountLessEqual:         recommendCountLessEqual,
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
