package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/go-bbs"
)

// TODO: go-bbs lacks following interfaces, remove when go-bbs will have implemented functions
// type PostsLimitedBoardRecord interface {
// 	EnableNewPost() bool
// }

// type postsLimitedBoardRecord struct{}

// func (r *postsLimitedBoardRecord) EnableNewPost() bool {
// 	return false
// }

func (repo *repository) GetBoards(_ context.Context) []bbs.BoardRecord {
	slog.Info("GetBoards", "boardRecords", len(repo.boardRecords))
	return repo.boardRecords
}

func (repo *repository) GetBoardArticle(_ context.Context, boardID, filename string) ([]byte, error) {
	return repo.db.ReadBoardArticleFile(boardID, filename)
}

func (repo *repository) GetBoardArticleRecords(_ context.Context, boardID string, offset, length uint) ([]bbs.ArticleRecord, error) {
	records, err := repo.db.ReadBoardArticleRecordsFile(boardID)
	if err != nil {
		return nil, err
	}
	return paginateArticleRecords(records, offset, length), nil
}

func paginateArticleRecords(records []bbs.ArticleRecord, offset, length uint) []bbs.ArticleRecord {
	if offset >= uint(len(records)) || length == 0 {
		return []bbs.ArticleRecord{}
	}

	start := int(offset)
	end := len(records)
	remaining := uint(end - start)
	if length < remaining {
		end = start + int(length)
	}

	return records[start:end]
}

func (repo *repository) GetBoardTreasureRecords(_ context.Context, boardID string, treasureIDs []string) ([]bbs.ArticleRecord, error) {
	return repo.db.ReadBoardTreasureRecordsFile(boardID, treasureIDs)
}

// func (repo *repository) GetBoardPostsLimit(_ context.Context, boardID string) (PostsLimitedBoardRecord, error) {
// 	// TODO: replace postsLimitedBoardRecord to real bbs record
// 	return &postsLimitedBoardRecord{}, nil
// }

func loadBoardFile(db *bbs.DB) ([]bbs.BoardRecord, error) {
	var logger = logging.NewLogger()
	boardRecords, err := db.ReadBoardRecords()
	if err != nil {
		logger.Errorf("get board header error: %v", err)
		return nil, fmt.Errorf("failed to read board records: %w", err)
	}
	for index, board := range boardRecords {
		logger.Debugf("loaded %d %v", index, board.BoardID())
	}
	return boardRecords, nil
}
