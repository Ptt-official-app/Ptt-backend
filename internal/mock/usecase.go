package mock

import (
	"context"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

// implements usecase.Usecase
type MockUsecase struct {
}

func NewMockUsecase() *MockUsecase {
	return &MockUsecase{}
}

func (u *MockUsecase) GetUserByID(ctx context.Context, userID string) (bbs.UserRecord, error) {
	panic("Not implemented")
}
func (u *MockUsecase) GetUserFavorites(ctx context.Context, userID string) ([]interface{}, error) {
	panic("Not implemented")
}
func (u *MockUsecase) GetUserInformation(ctx context.Context, userID string) (map[string]interface{}, error) {
	panic("Not implemented")
}

func (u *MockUsecase) GetBoardByID(ctx context.Context, userID string) (bbs.UserRecord, error) {
	panic("Not implemented")
}
func (u *MockUsecase) GetBoards(ctx context.Context, userID string) []bbs.BoardRecord {
	panic("Not implemented")
}
func (u *MockUsecase) GetClasses(ctx context.Context, userID, classID string) []bbs.BoardRecord {
	panic("Not implemented")
}
func (u *MockUsecase) GetBoardArticles(ctx context.Context, boardID string) []interface{} {
	panic("Not implemented")
}
func (u *MockUsecase) GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error) {
	panic("Not implemented")
}
func (u *MockUsecase) GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{} {
	panic("Not implemented")
}

func (u *MockUsecase) CreateAccessTokenWithUsername(username string) string {
	panic("Not implemented")
}
func (u *MockUsecase) GetUserIdFromToken(token string) (string, error) {
	panic("Not implemented")
}
func (u *MockUsecase) CheckPermission(token string, permissionId []usecase.Permission, userInfo map[string]string) error {
	panic("Not implemented")
}
