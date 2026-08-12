package usecase

// implements Repository
type MockRepository struct{}

// To make sure MockRepository implements the consumer-owned Repository.
var _ Repository = &MockRepository{}
