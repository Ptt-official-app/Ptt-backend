package usecase

import (
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

// implements repository.Repository
type MockRepository struct{}

// To make sure MockRepository implement repository.Repository
var _ repository.Repository = &MockRepository{}
