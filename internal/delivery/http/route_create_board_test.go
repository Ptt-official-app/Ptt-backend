package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Ptt-official-app/go-bbs"
)

func (usecase *MockUsecase) CreateBoard(_ context.Context, boardID, title string) (bbs.BoardRecord, error) {
	return NewMockBoardRecord("", boardID, title, false), nil
}

func TestCreateBoardResponse(t *testing.T) {
	mockUsecase := &MockUsecase{}
	delivery := NewHTTPDelivery(mockUsecase)

	form := url.Values{}
	form.Set("board_id", "testboard01")
	form.Set("title", "測試看板")
	req := httptest.NewRequest(http.MethodPost, "/v1/boards", strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "bearer token")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	delivery.routeBoards(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusCreated, rr.Body.String())
	}

	var response struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Data["id"] != "testboard01" || response.Data["title"] != "測試看板" {
		t.Fatalf("unexpected response data: %v", response.Data)
	}
}
