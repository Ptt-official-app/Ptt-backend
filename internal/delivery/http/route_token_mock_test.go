package http

import (
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (usecase *MockUsecase) CreateAccessTokenWithUsername(username string) string {
	return "token"
}

func (usecase *MockUsecase) GetUserIdFromToken(token string) (string, error) {
	return "id", nil
}

func (usecase *MockUsecase) CheckPermission(token string, permissionId []usecase.Permission, userInfo map[string]string) error {
	return nil
}
