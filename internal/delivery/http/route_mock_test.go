package http

import (
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

// implements usecase.Usecase
type MockUsecase struct {
}

func NewMockUsecase() usecase.Usecase {
	return &MockUsecase{}
}
