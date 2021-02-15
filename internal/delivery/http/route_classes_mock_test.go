package http

import (
	"context"
	"github.com/PichuChen/go-bbs"
)

type MockClassBoardRecord struct {
	boardID string
	isClass bool
	classID string
}

func NewMockClassBoardRecord(classID, boardID string, isClass bool) *MockClassBoardRecord {
	return &MockClassBoardRecord{boardID: boardID, isClass: isClass, classID: classID}
}

func (b *MockClassBoardRecord) BoardId() string { return b.boardID }

func (b *MockClassBoardRecord) Title() string { return "" }

func (b *MockClassBoardRecord) IsClass() bool { return b.isClass }

// ClassId should return the class id to which this board/class belongs.
func (b *MockClassBoardRecord) ClassId() string { return b.classID }

func (b *MockClassBoardRecord) BM() []string { return make([]string, 0) }

func (usecase *MockUsecase) GetClasses(ctx context.Context, userID, classID string) []bbs.BoardRecord {
	result := make([]bbs.BoardRecord, 0)
	result = append(result, NewMockClassBoardRecord(classID, "", true))
	result = append(result, NewMockClassBoardRecord(classID, "TEST", false))
	return result
}
