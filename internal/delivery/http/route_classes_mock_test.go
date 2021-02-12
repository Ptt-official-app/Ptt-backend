package http

import (
	"context"
	"github.com/PichuChen/go-bbs"
)

type MockClassBoardRecord struct {
	boardId string
	isClass bool
	classId string
}

func NewMockClassBoardRecord(classID, boardId string, isClass bool) *MockClassBoardRecord {
	return &MockClassBoardRecord{boardId: boardId, isClass: isClass, classId: classID}
}

func (b *MockClassBoardRecord) BoardId() string { return b.boardId }

func (b *MockClassBoardRecord) Title() string { return "" }

func (b *MockClassBoardRecord) IsClass() bool { return b.isClass }

// ClassId should return the class id to which this board/class belongs.
func (b *MockClassBoardRecord) ClassId() string { return b.classId }

func (b *MockClassBoardRecord) BM() []string { return make([]string, 0) }

func (usecase *MockUsecase) GetClasses(ctx context.Context, userID, classID string) []bbs.BoardRecord {
	result := make([]bbs.BoardRecord, 0)
	result = append(result, NewMockClassBoardRecord(classID, "", true))
	result = append(result, NewMockClassBoardRecord(classID, "TEST", false))
	return result
}
