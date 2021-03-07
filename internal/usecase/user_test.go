package usecase

import (
	"context"
	"testing"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
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

	if rec.UserId() != "pichu" {
		t.Errorf("getUserByID with pichu excepted userid: pichu, got %v", rec.UserId())
		return
	}

}

func TestGetUserArticles(t *testing.T) {

	userID := "userID"
	boardIDs := []string{"softjob", "techjob"}
	var mockRepository repository.Repository // FIXME: use concrete mock rather than mockRepository
	mockUsecase := NewUsecase(&config.Config{}, mockRepository)

	dataItems, err := mockUsecase.GetUserArticles(context.TODO(), boardIDs, userID)

	if err == nil {
		t.Errorf("getUserByID with userID excepted not nil error, got nil")
		return
	}

	if dataItems != nil {
		t.Errorf("getUserByID with userID excepted nil, got %v", dataItems)
		return
	}

}
