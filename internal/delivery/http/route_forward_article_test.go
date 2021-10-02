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
	"strings"
	"testing"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

// TestForwardArticleBadRequest test request post `/v1/boards/{}/articles/{}` post
// with no body
func TestForwardArticleBadRequest(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	req, err := http.NewRequest("POST", "/v1/boards/test/articles/test", nil)
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

// TestForwardArticleBadRequest test request post `/v1/boards/{}/articles/{}` post
// without valid email and board
func TestForwardArticleInternalServerError(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	v := url.Values{}
	v.Set("action", "forward_article")
	v.Set("board_id", "")
	v.Set("email", "")
	t.Logf("testing body: %v", v)
	req, err := http.NewRequest("POST", "/v1/boards/test/articles/test", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/", delivery.routeBoards)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestForwardArticleResponse(t *testing.T) {
	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	v := url.Values{}
	v.Set("action", "forward_article")
	v.Set("board_id", "SYSOP")
	v.Set("email", "pichu@tih.tw")
	t.Logf("testing body: %v", v)
	req, err := http.NewRequest("POST", "/v1/boards/test/articles/test", strings.NewReader(v.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/boards/", delivery.routeBoards)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

type MockForwardArticleToBoardUseCase struct {
	MockUsecase
	token                     string
	checkPermissionParameters func(permission usecase.Permission, permissionParameters map[string]string)
	forwardArticleToBoard     func(userID, boardID, filename, toBoard string) error
	forwardArticleToEmail     func(userID, boardID, filename, email string) error
}

func (mockUsecase *MockForwardArticleToBoardUseCase) CheckPermission(token string, permissionList []usecase.Permission, params map[string]string) error {
	for _, permission := range permissionList {
		switch permission {
		case usecase.PermissionForwardArticleToBoard:
			fallthrough
		case usecase.PermissionForwardArticleToEmail:
			if mockUsecase.checkPermissionParameters != nil {
				mockUsecase.checkPermissionParameters(permission, params)
			}
			if mockUsecase.token != token {
				return errors.New("token not match")
			}
		default:
			return errors.New("unexpected permission check: " + string(permission))
		}
	}
	return nil
}

func (mockUsecase *MockForwardArticleToBoardUseCase) ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (repository.ForwardArticleToBoardRecord, error) {
	if mockUsecase.forwardArticleToBoard != nil {
		return nil, mockUsecase.forwardArticleToBoard(userID, boardID, filename, boardName)
	}
	return nil, nil
}

func (mockUsecase *MockForwardArticleToBoardUseCase) ForwardArticleToEmail(ctx context.Context, userID, boardID, filename, email string) error {
	if mockUsecase.forwardArticleToEmail != nil {
		return mockUsecase.forwardArticleToEmail(userID, boardID, filename, email)
	}
	return nil
}

func (usecase *MockForwardArticleToBoardUseCase) GetUserIDFromToken(token string) (string, error) {
	if token == "" {
		return "", errors.New("token not found")
	}
	return "id", nil
}

func TestForwardArticleFunction(t *testing.T) {
	mockUsecase := &MockForwardArticleToBoardUseCase{}
	delivery := NewHTTPDelivery(mockUsecase)
	boardID := "test"
	toBoardID := "target"
	email := "example@example.com"
	filename := "random-filename"
	userData, err := mockUsecase.GetUserByID(context.Background(), "id")
	if err != nil {
		t.Fatalf("get error \"%s\" when create mock user %s", err.Error(), "id")
	}
	token := mockUsecase.CreateAccessTokenWithUsername(userData.UserID())

	// check token error
	body := url.Values{}
	body.Add("action", "forward_article")
	body.Add("board_id", toBoardID)
	req := httptest.NewRequest("POST", fmt.Sprintf("/v1/boards/%s/articles/%s", boardID, filename), strings.NewReader(body.Encode()))
	rr := httptest.NewRecorder()
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	delivery.forwardArticle(rr, req, boardID, filename)
	expected := map[string]interface{}{
		"error":             "token is not pass",
		"error_description": "token not found",
	}
	actual := make(map[string]interface{})
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	if err != nil {
		t.Fatalf("got error: \"%s\" when parsing response body", err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v but got %v", expected, actual)
	}

	// check parameters error
	req = httptest.NewRequest("POST", fmt.Sprintf("/v1/boards/%s/articles/%s", boardID, filename), nil)
	rr = httptest.NewRecorder()
	req.Header.Add("Authorization", "bearer "+token)
	delivery.forwardArticle(rr, req, boardID, filename)
	if !reflect.DeepEqual(rr.Body.Bytes(), NewNoRequiredParameterError(req, "email or board_id")) {
		t.Fatalf("expected %s but got %s", string(NewNoRequiredParameterError(req, "email or board_id")), rr.Body.String())
	}

	// check forward article to board
	body = url.Values{}
	body.Add("action", "forward_article")
	body.Add("board_id", toBoardID)
	req = httptest.NewRequest("POST",
		fmt.Sprintf("/v1/boards/%s/articles/%s", boardID, filename),
		strings.NewReader(body.Encode()))
	rr = httptest.NewRecorder()
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	mockUsecase.token = token
	permissionChecked := false
	mockUsecase.checkPermissionParameters = func(permission usecase.Permission, permissionParameters map[string]string) {
		permissionChecked = true
		if permission != usecase.PermissionForwardArticleToBoard {
			t.Fatalf("expect permission check %s, but got %s", usecase.PermissionForwardArticleToBoard, permission)
		}
		if pBoardID, ok := permissionParameters["board_id"]; !ok || pBoardID != boardID {
			t.Fatalf("expect board id %s, but got %s", boardID, pBoardID)
		}
		if pToBoard, ok := permissionParameters["to_board"]; !ok || pToBoard != toBoardID {
			t.Fatalf("expect to board id %s, but got %s", toBoardID, pToBoard)
		}
		if pArticle, ok := permissionParameters["article_id"]; !ok || pArticle != filename {
			t.Fatalf("expect article id %s, but got %s", filename, pArticle)
		}
	}
	isForwardToBoard := false
	mockUsecase.forwardArticleToBoard = func(pUserID, pBoardID, pFilename, pToBoard string) error {
		isForwardToBoard = true
		if pUserID != userData.UserID() {
			t.Fatalf("expect user id %s, but got %s", userData.UserID(), pUserID)
		}
		if pBoardID != boardID {
			t.Fatalf("expected board id %s, but got %s", boardID, pBoardID)
		}
		if pFilename != filename {
			t.Fatalf("expect filename %s, but got %s", filename, pFilename)
		}
		if pToBoard != toBoardID {
			t.Fatalf("expect to board id %s, but got %s", toBoardID, pToBoard)
		}
		return nil
	}
	delivery.forwardArticle(rr, req, boardID, filename)
	if !permissionChecked {
		t.Fatalf("Permission Not Checked")
	}
	if !isForwardToBoard {
		t.Fatalf("forwardArticleToBoard function is not called")
	}

	// check forward article to email
	body = url.Values{}
	body.Add("action", "forward_article")
	body.Add("email", email)
	req = httptest.NewRequest("POST",
		fmt.Sprintf("/v1/boards/%s/articles/%s", boardID, filename),
		strings.NewReader(body.Encode()))
	rr = httptest.NewRecorder()
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	mockUsecase.token = token
	permissionChecked = false
	mockUsecase.checkPermissionParameters = func(permission usecase.Permission, permissionParameters map[string]string) {
		permissionChecked = true
		if permission != usecase.PermissionForwardArticleToEmail {
			t.Fatalf("expect permission check %s, but got %s", usecase.PermissionForwardArticleToEmail, permission)
		}
		if pBoardID, ok := permissionParameters["board_id"]; !ok || pBoardID != boardID {
			t.Fatalf("expect board id %s, but got %s", boardID, pBoardID)
		}
		if pToEmail, ok := permissionParameters["to_email"]; !ok || pToEmail != email {
			t.Fatalf("expect to board id %s, but got %s", toBoardID, pToEmail)
		}
		if pArticle, ok := permissionParameters["article_id"]; !ok || pArticle != filename {
			t.Fatalf("expect article id %s, but got %s", filename, pArticle)
		}
	}
	isForwardToEmail := false
	mockUsecase.forwardArticleToEmail = func(pUserID, pBoardID, pFilename, pToEmail string) error {
		isForwardToEmail = true
		if pUserID != userData.UserID() {
			t.Fatalf("expect user id %s, but got %s", userData.UserID(), pUserID)
		}
		if pBoardID != boardID {
			t.Fatalf("expected board id %s, but got %s", boardID, pBoardID)
		}
		if pFilename != filename {
			t.Fatalf("expect filename %s, but got %s", filename, pFilename)
		}
		if pToEmail != email {
			t.Fatalf("expect to board id %s, but got %s", email, pToEmail)
		}
		return nil
	}
	delivery.forwardArticle(rr, req, boardID, filename)
	if !permissionChecked {
		t.Fatalf("Permission Not Checked")
	}
	if !isForwardToEmail {
		t.Fatalf("forwardArticleToBoard function is not called")
	}

	// test forward to board failed
	body = url.Values{}
	body.Add("action", "forward_article")
	body.Add("board_id", toBoardID)
	req = httptest.NewRequest("POST",
		fmt.Sprintf("/v1/boards/%s/articles/%s", boardID, filename),
		strings.NewReader(body.Encode()))
	rr = httptest.NewRecorder()
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	mockUsecase.checkPermissionParameters = func(permission usecase.Permission, permissionParameters map[string]string) {
	}
	mockUsecase.forwardArticleToBoard = func(pUserID, pBoardID, pFilename, pToBoard string) error {
		return errors.New("forward failed")
	}
	delivery.forwardArticle(rr, req, boardID, filename)
	decoder := json.NewDecoder(rr.Body)
	actual = make(map[string]interface{})
	expected = map[string]interface{}{
		"error":             "no_permission_for_create_board_articles",
		"error_description": "user don't have permission to create an article in board " + toBoardID,
	}
	err = decoder.Decode(&actual)
	if err != nil {
		t.Fatalf("decode response body failed: %s", err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %s, but got %s", expected, actual)
	}

	// test forward to email failed
	body = url.Values{}
	body.Add("action", "forward_article")
	body.Add("email", email)
	req = httptest.NewRequest("POST",
		fmt.Sprintf("/v1/boards/%s/articles/%s", boardID, filename),
		strings.NewReader(body.Encode()))
	rr = httptest.NewRecorder()
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	mockUsecase.forwardArticleToEmail = func(pUserID, pBoardID, pFilename, pToBoard string) error {
		return errors.New("forward failed")
	}
	delivery.forwardArticle(rr, req, boardID, filename)
	decoder = json.NewDecoder(rr.Body)
	actual = make(map[string]interface{})
	expected = map[string]interface{}{
		"error":             "no_permission_for_forward_board_article_to_email",
		"error_description": "user don't have permission to forward article random-filename to email " + email,
	}
	err = decoder.Decode(&actual)
	if err != nil {
		t.Fatalf("decode response body failed: %s", err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %s, but got %s", expected, actual)
	}

	// test forward to board permission denied
	mockUsecase.token = ""
	body = url.Values{}
	body.Add("action", "forward_article")
	body.Add("board_id", toBoardID)
	req = httptest.NewRequest("POST",
		fmt.Sprintf("/v1/boards/%s/articles/%s", boardID, filename),
		strings.NewReader(body.Encode()))
	rr = httptest.NewRecorder()
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	delivery.forwardArticle(rr, req, boardID, filename)
	decoder = json.NewDecoder(rr.Body)
	actual = make(map[string]interface{})
	expected = map[string]interface{}{
		"error":             "no_permission_for_create_board_articles",
		"error_description": "user don't have permission to create an article in board " + toBoardID,
	}
	err = decoder.Decode(&actual)
	if err != nil {
		t.Fatalf("decode response body failed: %s", err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %s, but got %s", expected, actual)
	}

	// test forward to email permission denined
	body = url.Values{}
	body.Add("action", "forward_article")
	body.Add("email", email)
	req = httptest.NewRequest("POST",
		fmt.Sprintf("/v1/boards/%s/articles/%s", boardID, filename),
		strings.NewReader(body.Encode()))
	rr = httptest.NewRecorder()
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	mockUsecase.forwardArticleToEmail = func(pUserID, pBoardID, pFilename, pToBoard string) error {
		return errors.New("forward failed")
	}
	delivery.forwardArticle(rr, req, boardID, filename)
	decoder = json.NewDecoder(rr.Body)
	actual = make(map[string]interface{})
	expected = map[string]interface{}{
		"error":             "no_permission_for_forward_board_article_to_email",
		"error_description": "user don't have permission to forward article random-filename to email " + email,
	}
	err = decoder.Decode(&actual)
	if err != nil {
		t.Fatalf("decode response body failed: %s", err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %s, but got %s", expected, actual)
	}

}
