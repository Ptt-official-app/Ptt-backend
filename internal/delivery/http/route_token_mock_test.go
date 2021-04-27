package http

import (
	"context"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (usecase *MockUsecase) CreateAccessTokenWithUsername(username string) string {
	return "token"
}

func (usecase *MockUsecase) GetUserIDFromToken(token string) (string, error) {
	return "id", nil
}

func (usecase *MockUsecase) CheckPermission(ctx context.Context, token string, permissionID []usecase.Permission, userInfo map[string]string) error {
	return nil
}
