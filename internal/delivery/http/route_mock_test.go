package http

// implements Usecase
type MockUsecase struct {
	loginUserID string
	loginIP     string
}

func NewMockUsecase() Usecase {
	return &MockUsecase{}
}

func (usecase *MockUsecase) RecordLogin(userID, ip string) {
	usecase.loginUserID = userID
	usecase.loginIP = ip
}

var _ Usecase = &MockUsecase{}
