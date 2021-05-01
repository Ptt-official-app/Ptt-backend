package usecase

import (
	"context"
	"testing"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
)

func TestGetUserByID(t *testing.T) {

	resp := &MockRepository{}

	usecase := NewUsecase(&config.Config{}, resp, logging.DefaultDummyLogger)

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
	mockUsecase := NewUsecase(&config.Config{}, mockRepository, logging.DefaultDummyLogger)

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
	expectBoardID := "SYSOP"
	mockRepository := &MockRepository{}
	mockUsecase := NewUsecase(&config.Config{}, mockRepository, logging.DefaultDummyLogger)

	dataItems, err := mockUsecase.GetUserComments(context.TODO(), userID)

	if err != nil {
		t.Errorf("GetUserComment with %s expect nil, got %v", userID, err)
	}

	item, ok := dataItems[0].(map[string]interface{})
	if !ok {
		t.Errorf("unexpect type of item")
	} else if item["board_id"] != expectBoardID {
		t.Errorf(`item["board_id"] expect %s, got %s`, expectBoardID, item["board_id"])
	}
}
