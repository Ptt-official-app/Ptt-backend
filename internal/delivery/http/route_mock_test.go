package http

import (
	"github.com/Ptt-official-app/Ptt-backend/internal/mail"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

// implements usecase.Usecase
type MockUsecase struct {
	loginUserID string
	loginIP     string
}

func NewMockUsecase() usecase.Usecase {
	return &MockUsecase{}
}

func (usecase *MockUsecase) RecordLogin(userID, ip string) {
	usecase.loginUserID = userID
	usecase.loginIP = ip
}

func (usecase *MockUsecase) UpdateMail(mail mail.Mail) error {
	return nil
}
