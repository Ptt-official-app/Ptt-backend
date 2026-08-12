package repository

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Ptt-official-app/go-bbs"
	_ "github.com/Ptt-official-app/go-bbs/pttbbs"
)

func TestCreateBoardPersistsAndUpdatesCache(t *testing.T) {
	home := t.TempDir()
	if err := os.WriteFile(filepath.Join(home, ".BRD"), nil, 0o600); err != nil {
		t.Fatalf("create empty .BRD: %v", err)
	}

	db, err := bbs.Open("pttbbs", home)
	if err != nil {
		t.Fatalf("open pttbbs: %v", err)
	}
	repo := &repository{db: db, boardRecords: []bbs.BoardRecord{}}

	created, err := repo.CreateBoard(context.Background(), "testboard01", "測試看板")
	if err != nil {
		t.Fatalf("CreateBoard() error = %v", err)
	}
	if created.BoardID() != "testboard01" || created.Title() != "測試看板" {
		t.Fatalf("created board = %q/%q", created.BoardID(), created.Title())
	}

	cached := repo.GetBoards(context.Background())
	if len(cached) != 1 || cached[0].BoardID() != "testboard01" {
		t.Fatalf("cached boards = %v", cached)
	}

	persisted, err := db.ReadBoardRecords()
	if err != nil {
		t.Fatalf("ReadBoardRecords() error = %v", err)
	}
	if len(persisted) != 1 || persisted[0].BoardID() != "testboard01" || persisted[0].Title() != "測試看板" {
		t.Fatalf("persisted boards = %v", persisted)
	}

	if _, err := repo.CreateBoard(context.Background(), "TESTBOARD01", "重複"); !errors.Is(err, ErrBoardExists) {
		t.Fatalf("duplicate CreateBoard() error = %v, want ErrBoardExists", err)
	}
}

func TestCreateBoardValidatesInput(t *testing.T) {
	repo := &repository{}
	for _, testCase := range []struct {
		boardID string
		title   string
	}{
		{boardID: "../escape", title: "bad"},
		{boardID: "board with space", title: "bad"},
		{boardID: "1234567890123", title: "too long"},
		{boardID: "valid", title: ""},
	} {
		if _, err := repo.CreateBoard(context.Background(), testCase.boardID, testCase.title); !errors.Is(err, ErrInvalidBoard) {
			t.Fatalf("CreateBoard(%q, %q) error = %v, want ErrInvalidBoard", testCase.boardID, testCase.title, err)
		}
	}
}
