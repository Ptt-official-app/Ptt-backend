package http

import (
	"context"
	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func (usecase *MockUsecase) GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error) {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetBoards(ctx context.Context, userID string) []bbs.BoardRecord {
	panic("Not implemented")
}

func (usecase *MockUsecase) GetPopularBoards(ctx context.Context) ([]bbs.BoardRecord, error) {
	result := make([]bbs.BoardRecord, 0)
	result = append(result, NewMockBoardRecord("SYSOP", "", "嘰哩 ◎站長好!", true))
	result = append(result, NewMockBoardRecord("junk", "TEST", "發電 ◎雜七雜八的垃圾", false))
	return result, nil
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

type MockBoardRecord struct {
	boardID string
	title   string
	isClass bool
	classID string
}

func NewMockBoardRecord(classID, boardID, title string, isClass bool) *MockBoardRecord {
	return &MockBoardRecord{boardID: boardID, title: title, isClass: isClass, classID: classID}
}

func (b *MockBoardRecord) BoardId() string { return b.boardID }
func (b *MockBoardRecord) Title() string { return b.title }
func (b *MockBoardRecord) IsClass() bool { return b.isClass }
func (b *MockBoardRecord) ClassId() string { return b.classID }
func (b *MockBoardRecord) BM() []string { return make([]string, 0) }
