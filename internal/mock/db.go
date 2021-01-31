package mock

import (
	"github.com/PichuChen/go-bbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/db"
)

// implements db.DB
type MockDB struct {
}

func NewMockDB() db.DB {
	return &MockDB{}
}

func (db *MockDB) ReadBoardRecords() ([]bbs.BoardRecord, error) {
	panic("Not implemented")
}
func (db *MockDB) ReadUserRecords() ([]bbs.UserRecord, error) {
	panic("Not implemented")
}
func (db *MockDB) ReadBoardArticleFile(boardId string, filename string) ([]byte, error) {
	panic("Not implemented")
}
func (db *MockDB) ReadBoardArticleRecordsFile(boardId string) ([]bbs.ArticleRecord, error) {
	panic("Not implemented")
}
func (db *MockDB) ReadBoardTreasureRecordsFile(boardId string, treasureId []string) ([]bbs.ArticleRecord, error) {
	panic("Not implemented")
}
func (db *MockDB) ReadUserFavoriteRecords(userId string) ([]bbs.FavoriteRecord, error) {
	panic("Not implemented")
}
