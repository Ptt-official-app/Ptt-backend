package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
)

func TestGetUserByID(t *testing.T) {
	resp := &MockRepository{}

	usecase := NewUsecase(&config.Config{}, resp)

	rec, err := usecase.GetUserByID(context.TODO(), "not-exist-user-id")
	if err == nil {
		t.Errorf("getUserByID with not-exist-user-id excepted not nil error, got nil")
		return
	}

	if rec != nil {
		t.Errorf("getUserByID with not-exist-user-id excepted nil, got %v", rec)
		return
	}

	rec, err = usecase.GetUserByID(context.TODO(), "pichu")
	if err != nil {
		t.Errorf("getUserByID with pichu excepted err == nil, got %v", err)
		return
	}

	if rec.UserID() != "pichu" {
		t.Errorf("getUserByID with pichu excepted userid: pichu, got %v", rec.UserID())
		return
	}
}

func TestGetUserInformation_InputNotExistUser_ReturnError(t *testing.T) {
	userID := "not-exist-user-id"
	errMsg := fmt.Sprintf("get userrec for %s failed", userID)
	repo := &MockRepository{}
	usecase := NewUsecase(&config.Config{}, repo)

	data, err := usecase.GetUserInformation(context.TODO(), userID)
	// TODO: check return error message with error object
	if err.Error() != errMsg {
		t.Errorf("GetUserInformation with %s expected error, got %v", userID, err)
	}

	if data != nil {
		t.Errorf("GetUserInformation with %s expect nil data, got %v", userID, data)
	}
}

func TestGetUserInformation_InputPichu_ReturnData(t *testing.T) {
	userID := "pichu"
	repo := &MockRepository{}
	usecase := NewUsecase(&config.Config{}, repo)

	data, err := usecase.GetUserInformation(context.TODO(), userID)
	if err != nil {
		t.Errorf("GetUserInformation with %s expected nil, got %v", userID, err)
	}

	if data == nil {
		t.Errorf("GetUserInformation with %s expected data, got nil", userID)
	}

	if data["user_id"] != userID {
		t.Errorf("GetUserInformation with %s expect %s, got %s", userID, userID, data["user_id"])
	}
}

func TestGetUserArticles(t *testing.T) {
	userID := "user"
	mockRepository := &MockRepository{}
	mockUsecase := NewUsecase(&config.Config{}, mockRepository)

	dataItems, err := mockUsecase.GetUserArticles(context.TODO(), userID)
	if err != nil {
		t.Errorf("GetUserArticles with userID excepted nil error, got %v", err)
		return
	}

	if dataItems == nil {
		t.Errorf("GetUserArticles with userID excepted not nil, got %v", dataItems)
		return
	}
}

func TestGetUserComments(t *testing.T) {
	userID := "user"
	mockRepository := &MockRepository{}
	mockUsecase := NewUsecase(&config.Config{}, mockRepository)

	_, err := mockUsecase.GetUserComments(context.TODO(), userID)
	if err != nil {
		t.Errorf("GetUserComment with %s expect nil, got %v", userID, err)
	}
}

func TestGetUserDrafts(t *testing.T) {
	userID := "user"
	mockRepository := &MockRepository{}
	mockUsecase := NewUsecase(&config.Config{}, mockRepository)

	// case 1: valid draftID
	actualValue, _ := mockUsecase.GetUserDrafts(context.TODO(), userID, "0")
	expectedValue := "this is a draft"
	if expectedValue != string(actualValue.Raw()) {
		t.Errorf("returned unexpected value: got %v want value %v",
			actualValue, expectedValue)
	}

	// case 2: invalid draftID
	_, err := mockUsecase.GetUserDrafts(context.TODO(), userID, "10")
	if err == nil {
		t.Error("returned unexpected error")
	}
}

func TestUpdateUserDraft(t *testing.T) {
	userID := "user"
	mockRepository := &MockRepository{}
	mockUsecase := NewUsecase(&config.Config{}, mockRepository)

	actualValue, _ := mockUsecase.UpdateUserDraft(context.TODO(), userID, "0", []byte("this is a draft"))
	expectedValue := "this is a draft"

	if expectedValue != string(actualValue.Raw()) {
		t.Errorf("returned unexpected value: got %v want value %v",
			actualValue, expectedValue)
	}
}

func TestDeleteUserDraft(t *testing.T) {
	userID := "user"
	mockRepository := &MockRepository{}
	mockUsecase := NewUsecase(&config.Config{}, mockRepository)

	// case 1: valid draftID
	err := mockUsecase.DeleteUserDraft(context.TODO(), userID, "0")
	if err != nil {
		t.Error("returned unexpected error")
	}

	// case 2: invalid draftID
	err = mockUsecase.DeleteUserDraft(context.TODO(), userID, "10")
	if err == nil {
		t.Error("returned unexpected error")
	}
}
