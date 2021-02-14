package http

import (
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"

	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// implements usecase.Usecase
type MockBoardsUsecase struct {
}

func NewMockBoardsUsecase() usecase.Usecase {
	return &MockBoardsUsecase{}
}

// usecase/user.go
func (usecase *MockBoardsUsecase) GetUserByID(ctx context.Context, userID string) (bbs.UserRecord, error) {
	panic("Not implemented")
}

func (usecase *MockBoardsUsecase) GetUserFavorites(ctx context.Context, userID string) ([]interface{}, error) {
	panic("Not implemented")
}

func (usecase *MockBoardsUsecase) GetUserInformation(ctx context.Context, userID string) (map[string]interface{}, error) {
	panic("Not implemented")
}

// usecase/board.go
func (usecase *MockBoardsUsecase) GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error) {
	panic("Not implemented")
}

func (usecase *MockBoardsUsecase) GetBoards(ctx context.Context, userID string) []bbs.BoardRecord {
	//panic("Not implemented")
	return []bbs.BoardRecord{}
}

func (usecase *MockBoardsUsecase) GetClasses(ctx context.Context, userID, classID string) []bbs.BoardRecord {
	panic("Not implemented")
}

func (usecase *MockBoardsUsecase) GetBoardArticles(ctx context.Context, boardID string) []interface{} {
	panic("Not implemented")
}

func (usecase *MockBoardsUsecase) GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error) {
	panic("Not implemented")
}

func (usecase *MockBoardsUsecase) GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{} {
	panic("Not implemented")
}

// usecase/token.go
func (usecase *MockBoardsUsecase) CreateAccessTokenWithUsername(username string) string {
	return "token"
}

func (usecase *MockBoardsUsecase) GetUserIdFromToken(token string) (string, error) {
	return "id", nil
}

func (usecase *MockBoardsUsecase) CheckPermission(token string, permissionId []usecase.Permission, userInfo map[string]string) error {
	return nil
}

func TestGetBoardList (t *testing.T) {
	userID := "id"
	usecase := NewMockBoardsUsecase()
	delivery := NewHTTPDelivery(usecase)

	req, err := http.NewRequest("GET", "/v1/boards", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	w := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("v1/boards", delivery.routeBoards)
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
