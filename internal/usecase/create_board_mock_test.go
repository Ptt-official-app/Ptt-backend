package usecase

import (
	"context"

	"github.com/Ptt-official-app/go-bbs"
)

func (repo *MockRepository) CreateBoard(_ context.Context, boardID, title string) (bbs.BoardRecord, error) {
	return NewMockBoardRecord("", boardID, title, false), nil
}
