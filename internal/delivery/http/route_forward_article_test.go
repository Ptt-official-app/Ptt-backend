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
	UseCase "github.com/Ptt-official-app/Ptt-backend/internal/usecase"
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
	checkPermissionParameters func(permission UseCase.Permission, permissionParameters map[string]string)
	forwardArticleToBoard     func(userID, boardID, filename, toBoard string) error
	forwardArticleToEmail     func(userID, boardID, filename, email string) error
}

func (usecase *MockForwardArticleToBoardUseCase) CheckPermission(token string, permissionList []UseCase.Permission, params map[string]string) error {
	for _, permission := range permissionList {
		switch permission {
		case UseCase.PermissionForwardArticleToBoard:
			fallthrough
		case UseCase.PermissionForwardArticleToEmail:
			usecase.checkPermissionParameters(permission, params)
			if usecase.token != token {
				return errors.New("token not match")
			}
		default:
			return errors.New("unexpected permission check: " + string(permission))
		}
	}
	return nil
}

func (usecase *MockForwardArticleToBoardUseCase) ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (repository.ForwardArticleToBoardRecord, error) {
	return nil, usecase.forwardArticleToBoard(userID, boardID, filename, boardName)
}

func (usecase *MockForwardArticleToBoardUseCase) ForwardArticleToEmail(ctx context.Context, userID, boardID, filename, email string) error {
	return usecase.forwardArticleToEmail(userID, boardID, filename, email)
}

func TestForwardArticleFunction(t *testing.T) {
	usecase := &MockForwardArticleToBoardUseCase{}
	delivery := NewHTTPDelivery(usecase)
	boardID := "test"
	toBoardID := "target"
	email := "example@example.com"
	filename := "random-filename"
	userData, err := usecase.GetUserByID(context.Background(), "id")
	if err != nil {
		t.Fatalf("get error \"%s\" when create mock user %s", err.Error(), "id")
	}
	token := usecase.CreateAccessTokenWithUsername(userData.UserID())

	// check token error
	req := httptest.NewRequest("POST", fmt.Sprintf("/v1/boards/%s/articles/%s", boardID, filename), nil)
	rr := httptest.NewRecorder()
	delivery.forwardArticle(rr, req, boardID, filename)
	expected := map[string]interface{}{
		"error":             "no_required_parameter",
		"error_description": "required parameter: email or board_id not found",
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
	body := url.Values{}
	body.Add("action", "forward_article")
	body.Add("board_id", toBoardID)
	req = httptest.NewRequest("POST",
		fmt.Sprintf("/v1/boards/%s/articles/%s", boardID, filename),
		strings.NewReader(body.Encode()))
	rr = httptest.NewRecorder()
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	usecase.token = token
	permissionChecked := false
	usecase.checkPermissionParameters = func(permission UseCase.Permission, permissionParameters map[string]string) {
		permissionChecked = true
		if permission != UseCase.PermissionForwardArticleToBoard {
			t.Fatalf("expect permission check %s, but got %s", UseCase.PermissionForwardArticleToBoard, permission)
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
	usecase.forwardArticleToBoard = func(pUserID, pBoardID, pFilename, pToBoard string) error {
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
	usecase.token = token
	permissionChecked = false
	usecase.checkPermissionParameters = func(permission UseCase.Permission, permissionParameters map[string]string) {
		permissionChecked = true
		if permission != UseCase.PermissionForwardArticleToEmail {
			t.Fatalf("expect permission check %s, but got %s", UseCase.PermissionForwardArticleToBoard, permission)
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
	usecase.forwardArticleToEmail = func(pUserID, pBoardID, pFilename, pToEmail string) error {
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
}
