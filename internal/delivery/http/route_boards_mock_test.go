package http

import (
	"context"
	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (usecase *MockUsecase) GetBoardByID(ctx context.Context, BoardID string) (bbs.BoardRecord, error) {
	BoardRecord := NewMockBoardRecord(BoardID, "")
	return BoardRecord, nil
}

func (usecase *MockUsecase) GetBoards(ctx context.Context, userID string) []bbs.BoardRecord {
	return []bbs.BoardRecord{}
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
	BoardID string
	ClassID string
}

func NewMockBoardRecord(BoardID string, ClassID string) *MockBoardRecord {
	return &MockBoardRecord{
		BoardID: BoardID,
		ClassID: ClassID,
	}
}

// BoardId will be replaced as BoardID in the future.
func (b *MockBoardRecord) BoardId() string { return b.BoardID }

func (b *MockBoardRecord) Title() string { return "" }

func (b *MockBoardRecord) IsClass() bool { return true }

// ClassId should return the class id to which this board/class belongs.
func (b *MockBoardRecord) ClassId() string { return b.ClassID }

func (b *MockBoardRecord) BM() []string { return []string{} }
