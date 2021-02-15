package http

import (
	"context"

	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (usecase *MockUsecase) GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error) {
	BoardRecord := NewMockBoardRecord(boardID)
	return BoardRecord, nil
}

func (usecase *MockUsecase) GetBoards(ctx context.Context, userID string) []bbs.BoardRecord {
	return []bbs.BoardRecord{}
}

func (usecase *MockUsecase) GetClasses(ctx context.Context, userID, classID string) []bbs.BoardRecord {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoardArticles(ctx context.Context, boardID string, cond *usecase.ArticleSearchCond) []interface{} {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error) {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{} {
	panic("Not implemented")
}

type MockBoardRecord struct{
	boardID string
}

func NewMockBoardRecord(boardID string) *MockBoardRecord {
	return &MockBoardRecord{
		boardID: boardID,
	}
}

func (b *MockBoardRecord) BoardId() string { return b.boardID }

func (b *MockBoardRecord) Title() string { return "" }

func (b *MockBoardRecord) IsClass() bool { return true }

// ClassId should return the class id to which this board/class belongs.
func (b *MockBoardRecord) ClassId() string { return "" }

func (b *MockBoardRecord) BM() []string { return []string{} }
