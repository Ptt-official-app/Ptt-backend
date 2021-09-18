package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
	UseCase "github.com/Ptt-official-app/Ptt-backend/internal/usecase"
	"github.com/Ptt-official-app/go-bbs"
)

func TestGetBoardArticlesBadRequest(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	title := ""
	author := ""
	recommendCountGe := "qwerty"
	v := url.Values{}
	v.Set("title", title)
	v.Set("author", author)
	v.Set("recommend_count_ge", recommendCountGe)
	uri := fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/", delivery.routeBoards)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestGetBoardArticlesResponse(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	titleContain := "test_posts"
	author := "test01"
	v := url.Values{}
	v.Set("title_contain", titleContain)
	v.Set("author", author)
	uri := fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/", delivery.routeBoards)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	actualResponseMap := map[string]interface{}{}
	if err := json.Unmarshal(rr.Body.Bytes(), &actualResponseMap); err != nil {
		t.Error(err.Error())
	}
}

type MockArticleUsecase struct {
	MockUsecase
	token            string
	getBoardArticles func(ctx context.Context, boardID string, cond *usecase.ArticleSearchCond) []bbs.ArticleRecord
}

func (usecase *MockArticleUsecase) GetBoardArticles(ctx context.Context, boardID string, cond *UseCase.ArticleSearchCond) []bbs.ArticleRecord {
	if usecase.getBoardArticles != nil {
		return usecase.getBoardArticles(ctx, boardID, cond)
	}
	return []bbs.ArticleRecord{}
}

func (usecase *MockArticleUsecase) CheckPermission(token string, permissionID []usecase.Permission, userInfo map[string]string) error {
	if token != usecase.token || token == "" {
		return errors.New("invalid token")
	}
	return nil
}

func TestGetBoardArticlesFunction(t *testing.T) {
	userID := "id"
	usecase := &MockArticleUsecase{}
	delivery := NewHTTPDelivery(usecase)
	boardID := "test"

	// permission denied test
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	delivery.getBoardArticles(rr, req, boardID)

	decoder := json.NewDecoder(rr.Body)
	actual := make(map[string]interface{})
	err := decoder.Decode(&actual)
	if err != nil {
		t.Fatalf("decode body error: %s", err.Error())
	}
	expect := map[string]interface{}{
		"error":             "no_permission_for_read_board_articles",
		"error_description": "user don't have permission for read board test",
	}
	if !reflect.DeepEqual(expect, actual) {
		t.Fatalf("expect %v, but get %v", expect, actual)
	}

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Unexpected status code: %d\n", rr.Code)
	}

	// title search test
	targetRecord := &MockArticleRecord{
		title:          "hello",
		owner:          "test01",
		filename:       "none",
		modified:       time.Now(),
		recommendCount: 12,
	}
	usecase.getBoardArticles = func(ctx context.Context, boardID string, cond *UseCase.ArticleSearchCond) []bbs.ArticleRecord {
		if cond.Title != targetRecord.title || cond.Author != targetRecord.owner {
			t.Fatalf("Title search failed, expected %s and %s, but got %s and %s", targetRecord.title, targetRecord.owner, cond.Title, cond.Author)
		}
		return []bbs.ArticleRecord{
			targetRecord,
		}
	}
	v := url.Values{}
	v.Set("title_contain", targetRecord.title)
	v.Set("author", targetRecord.owner)
	uri := fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req = httptest.NewRequest("GET", uri, nil)
	rr = httptest.NewRecorder()
	token := usecase.CreateAccessTokenWithUsername(userID)
	usecase.token = token
	req.Header.Add("Authorization", "bearer "+token)
	delivery.getBoardArticles(rr, req, "test")
	actual = make(map[string]interface{})
	decoder = json.NewDecoder(rr.Body)
	err = decoder.Decode(&actual)
	if err != nil {
		t.Fatalf("decode body error: %s", err.Error())
	}
	expect = make(map[string]interface{})
	fmt.Println(targetRecord.Date())
	_ = json.Unmarshal([]byte(fmt.Sprintf(`
	{
		"data": {
			"items": [{
				"filename": "none",
				"modified_time": "%s",
				"money": 0,
				"owner": "test01",
				"post_date": "",
				"recommend_count": 12,
				"title": "hello"
			}]
		}
	}`, targetRecord.Modified().Format(time.RFC3339Nano))), &expect)
	if !reflect.DeepEqual(expect, actual) {
		t.Fatalf("expect %v, but get %v", expect, actual)
	}

	// test recommend search
	// search 20 < recommend < 40, aka 21 <= recommend <= 39.
	searchOption := &UseCase.ArticleSearchCond{
		RecommendCountGreaterEqual:      21,
		RecommendCountLessEqual:         39,
		RecommendCountGreaterEqualIsSet: true,
		RecommendCountLessEqualIsSet:    true,
	}
	usecase.getBoardArticles = func(ctx context.Context, boardID string, cond *UseCase.ArticleSearchCond) []bbs.ArticleRecord {
		if !reflect.DeepEqual(searchOption, cond) {
			t.Fatalf("expect search option %v, but got %v", searchOption, cond)
		}
		return []bbs.ArticleRecord{}
	}
	// the range should be 20 < recommend < 40
	v = url.Values{}
	v.Set("recommend_count_ge", "10")
	v.Set("recommend_count_gt", "20")

	v.Set("recommend_count_le", "50")
	v.Set("recommend_count_lt", "40")

	uri = fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "bearer "+token)
	rr = httptest.NewRecorder()
	delivery.getBoardArticles(rr, req, "test")

	// the range should be 20 < recommend < 40
	v = url.Values{}
	v.Set("recommend_count_ge", "21")
	v.Set("recommend_count_gt", "0")

	v.Set("recommend_count_le", "39")
	v.Set("recommend_count_lt", "104")
	uri = fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "bearer "+token)
	rr = httptest.NewRecorder()
	delivery.getBoardArticles(rr, req, "test")

	v = url.Values{}
	v.Set("recommend_count_ge", "21")
	v.Set("recommend_count_le", "39")
	uri = fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "bearer "+token)
	rr = httptest.NewRecorder()
	delivery.getBoardArticles(rr, req, "test")

	v = url.Values{}
	v.Set("recommend_count_gt", "20")
	v.Set("recommend_count_lt", "40")
	uri = fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "bearer "+token)
	rr = httptest.NewRecorder()
	delivery.getBoardArticles(rr, req, "test")

	searchOption.RecommendCountGreaterEqualIsSet = false
	searchOption.RecommendCountGreaterEqual = 0
	v = url.Values{}
	v.Set("recommend_count_lt", "40")
	uri = fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "bearer "+token)
	rr = httptest.NewRecorder()
	delivery.getBoardArticles(rr, req, "test")

	searchOption.RecommendCountGreaterEqualIsSet = true
	searchOption.RecommendCountGreaterEqual = 21
	searchOption.RecommendCountLessEqualIsSet = false
	searchOption.RecommendCountLessEqual = 0
	v = url.Values{}
	v.Set("recommend_count_gt", "20")
	uri = fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "bearer "+token)
	rr = httptest.NewRecorder()
	delivery.getBoardArticles(rr, req, "test")

	// test recommend is not integer
	usecase.getBoardArticles = func(ctx context.Context, boardID string, cond *UseCase.ArticleSearchCond) []bbs.ArticleRecord {
		return []bbs.ArticleRecord{}
	}

	v = url.Values{}
	v.Set("recommend_count_gt", "a")
	uri = fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "bearer "+token)
	rr = httptest.NewRecorder()
	delivery.getBoardArticles(rr, req, "test")
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expect status code %d but got %d", http.StatusBadRequest, rr.Code)
	}

	v = url.Values{}
	v.Set("recommend_count_ge", "a")
	uri = fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "bearer "+token)
	rr = httptest.NewRecorder()
	delivery.getBoardArticles(rr, req, "test")
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expect status code %d but got %d", http.StatusBadRequest, rr.Code)
	}

	v = url.Values{}
	v.Set("recommend_count_lt", "a")
	uri = fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "bearer "+token)
	rr = httptest.NewRecorder()
	delivery.getBoardArticles(rr, req, "test")
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expect status code %d but got %d", http.StatusBadRequest, rr.Code)
	}

	v = url.Values{}
	v.Set("recommend_count_le", "a")
	uri = fmt.Sprintf("/v1/boards/test/articles?%s", v.Encode())
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "bearer "+token)
	rr = httptest.NewRecorder()
	delivery.getBoardArticles(rr, req, "test")
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expect status code %d but got %d", http.StatusBadRequest, rr.Code)
	}
}
