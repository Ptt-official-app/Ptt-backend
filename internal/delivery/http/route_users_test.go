package http

import (
	"github.com/PichuChen/go-bbs"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"

	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockUserRecord struct {
	userID string
}

func NewMockUserRecord(userID string) *MockUserRecord { return &MockUserRecord{userID: userID} }
func (u *MockUserRecord) UserID() string              { return u.userID }

// HashedPassword return user hashed password, it only for debug,
// If you want to check is user password correct, please use
// VerifyPassword insteaded.
func (u *MockUserRecord) HashedPassword() string { return "" }

// VerifyPassword will check user's password is OK. it will return null
// when OK and error when there are something wrong
func (u *MockUserRecord) VerifyPassword(password string) error { return nil }

// Nickname return a string for user's nickname, this string may change
// depend on user's mood, return empty string if this bbs system do not support
func (u *MockUserRecord) Nickname() string { return "" }

// RealName return a string for user's real name, this string may not be changed
// return empty string if this bbs system do not support
func (u *MockUserRecord) RealName() string { return "" }

// NumLoginDays return how many days this have been login since account created.
func (u *MockUserRecord) NumLoginDays() int { return 0 }

// NumPosts return how many posts this user has posted.
func (u *MockUserRecord) NumPosts() int { return 0 }

// Money return the money this user have.
func (u *MockUserRecord) Money() int { return 0 }

// LastLogin return last login time of user
func (u *MockUserRecord) LastLogin() time.Time { return time.Now() }

// LastHost return last login host of user, it is IPv4 address usually, but it
// could be domain name or IPv6 address.
func (u *MockUserRecord) LastHost() string { return "" }

// implements usecase.Usecase
type MockUsecase struct {
}

func NewMockUsecase() usecase.Usecase {
	return &MockUsecase{}
}

// usecase/user.go
func (usecase *MockUsecase) GetUserByID(ctx context.Context, userID string) (bbs.UserRecord, error) {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetUserFavorites(ctx context.Context, userID string) ([]interface{}, error) {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetUserInformation(ctx context.Context, userID string) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"user_id": "id",
	}
	return result, nil
}

// usecase/board.go
func (usecase *MockUsecase) GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error) {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoards(ctx context.Context, userID string) []bbs.BoardRecord {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetClasses(ctx context.Context, userID, classID string) []bbs.BoardRecord {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoardArticles(ctx context.Context, boardID string) []interface{} {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error) {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{} {
	panic("Not implemented")
}

// usecase/token.go
func (usecase *MockUsecase) CreateAccessTokenWithUsername(username string) string {
	return "token"
}

func (usecase *MockUsecase) GetUserIdFromToken(token string) (string, error) {
	panic("Not implemented")
}

func (usecase *MockUsecase) CheckPermission(token string, permissionId []usecase.Permission, userInfo map[string]string) error {
	return nil
}

func TestGetUserInformation(t *testing.T) {

	userID := "id"
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)
	req, err := http.NewRequest("GET", "/v1/users/SYSOP/information", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := usecase.CreateAccessTokenWithUsername(userID)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/users/", delivery.routeUsers)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responsedMap := map[string]interface{}{}
	json.Unmarshal(rr.Body.Bytes(), &responsedMap)
	t.Logf("got response %v", rr.Body.String())
	responsedData := responsedMap["data"].(map[string]interface{})
	if responsedData["user_id"] != userID {
		t.Errorf("handler returned unexpected body, user_id not match: got %v want userID %v",
			rr.Body.String(), userID)
	}
}

func TestParseUserPath(t *testing.T) {

	type TestCase struct {
		input          string
		expectedUserID string
		expectedItem   string
	}

	cases := []TestCase{
		{
			input:          "/v1/users/Pichu/information",
			expectedUserID: "Pichu",
			expectedItem:   "information",
		},
		{
			input:          "/v1/users/Pichu/",
			expectedUserID: "Pichu",
			expectedItem:   "",
		},
		{
			input:          "/v1/users/Pichu",
			expectedUserID: "Pichu",
			expectedItem:   "",
		},
	}

	for index, c := range cases {
		input := c.input
		expectedUserID := c.expectedUserID
		expectedItem := c.expectedItem
		actualUserID, actualItem, err := parseUserPath(input)
		if err != nil {
			t.Errorf("error on index %d, got: %v", index, err)

		}

		if actualUserID != expectedUserID {
			t.Errorf("userID not match on index %d, expected: %v, got: %v", index, expectedUserID, actualUserID)
		}

		if actualItem != expectedItem {
			t.Errorf("item not match on index %d, expected: %v, got: %v", index, expectedItem, actualItem)
		}

	}

}
