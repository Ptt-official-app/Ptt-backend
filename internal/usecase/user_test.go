package usecase

import (
	"context"
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
