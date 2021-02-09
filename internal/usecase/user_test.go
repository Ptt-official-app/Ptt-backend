package usecase

import (
	"github.com/Ptt-official-app/Ptt-backend/internal/config"

	"context"
	"testing"
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
