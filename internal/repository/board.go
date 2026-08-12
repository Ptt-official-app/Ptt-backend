package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/go-bbs"
	"github.com/Ptt-official-app/go-bbs/pttbbs"
)

var (
	ErrBoardExists  = errors.New("board already exists")
	ErrInvalidBoard = errors.New("invalid board")
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
	repo.boardMu.RLock()
	defer repo.boardMu.RUnlock()

	slog.Info("GetBoards", "boardRecords", len(repo.boardRecords))
	result := make([]bbs.BoardRecord, len(repo.boardRecords))
	copy(result, repo.boardRecords)
	return result
}

func validBoardID(boardID string) bool {
	if len(boardID) == 0 || len(boardID) > pttbbs.IDLength {
		return false
	}
	for _, r := range boardID {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			continue
		}
		return false
	}
	return true
}

func (repo *repository) CreateBoard(_ context.Context, boardID, title string) (bbs.BoardRecord, error) {
	boardID = strings.TrimSpace(boardID)
	title = strings.TrimSpace(title)
	if !validBoardID(boardID) || title == "" || len(bbs.Utf8ToBig5(title)) > pttbbs.BoardTitleLength {
		return nil, ErrInvalidBoard
	}

	repo.boardMu.Lock()
	defer repo.boardMu.Unlock()

	for _, board := range repo.boardRecords {
		if strings.EqualFold(board.BoardID(), boardID) {
			return nil, ErrBoardExists
		}
	}

	record, err := repo.db.NewBoardRecord(map[string]interface{}{
		"board_id": boardID,
		"title":    title,
	})
	if err != nil {
		return nil, fmt.Errorf("create board record: %w", err)
	}
	if err := repo.db.AddBoardRecord(record); err != nil {
		return nil, fmt.Errorf("persist board record: %w", err)
	}

	repo.boardRecords = append(repo.boardRecords, record)
	return record, nil
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
