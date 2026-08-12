package repository

import (
	"fmt"
	"sync"

	"github.com/Ptt-official-app/go-bbs"
)

type repository struct {
	db              *bbs.DB
	userRecords     []bbs.UserRecord
	boardMu         sync.RWMutex
	boardRecords    []bbs.BoardRecord
	popularArticles popularArticlesCache
}

func NewRepository(db *bbs.DB) (*repository, error) {
	userRecords, err := loadUserRecords(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load user records: %w", err)
	}

	boardRecords, err := loadBoardFile(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load board file: %w", err)
	}

	return &repository{
		db:           db,
		boardRecords: boardRecords,
		userRecords:  userRecords,
	}, nil
}
