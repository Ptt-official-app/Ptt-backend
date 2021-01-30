package repository

import (
	"context"
	"fmt"

	"github.com/PichuChen/go-bbs"
)

type boardRepository struct {
	db           *bbs.DB
	boardRecords []bbs.BoardRecord
}

func NewBoardRepository(db *bbs.DB) (BoardRepository, error) {
	boardRecords, err := loadBoardFile(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load board file: %w", err)
	}

	return &boardRepository{
		db:           db,
		boardRecords: boardRecords,
	}, nil
}

func (b *boardRepository) GetBoards(_ context.Context) []bbs.BoardRecord {
	return b.boardRecords
}

func (b *boardRepository) GetBoardArticle(_ context.Context, boardID, filename string) ([]byte, error) {
	return b.db.ReadBoardArticleFile(boardID, filename)
}

func (b *boardRepository) GetBoardArticleRecords(_ context.Context, boardID string) ([]bbs.ArticleRecord, error) {
	return b.db.ReadBoardArticleRecordsFile(boardID)
}

func (b *boardRepository) GetBoardTreasureRecords(_ context.Context, boardID string, treasureIDs []string) ([]bbs.ArticleRecord, error) {
	return b.db.ReadBoardTreasureRecordsFile(boardID, treasureIDs)
}

func loadBoardFile(db *bbs.DB) ([]bbs.BoardRecord, error) {
	boardRecords, err := db.ReadBoardRecords()
	if err != nil {
		logger.Errorf("get board header error: %v", err)
		return nil, fmt.Errorf("failed to read board records: %w", err)
	}
	for index, board := range boardRecords {
		logger.Debugf("loaded %d %v", index, board.BoardId())
	}
	return boardRecords, nil
}
