package http

import (
	"github.com/Ptt-official-app/Ptt-backend/internal/mail"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

// implements usecase.Usecase
type MockUsecase struct {
}

func NewMockUsecase() usecase.Usecase {
	return &MockUsecase{}
}

func (usecase *MockUsecase) UpdateMail(mail mail.Mail) error {
	return nil
}
