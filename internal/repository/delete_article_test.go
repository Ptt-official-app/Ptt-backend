package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Ptt-official-app/go-bbs"
	"github.com/Ptt-official-app/go-bbs/pttbbs"
)

func TestDeleteBoardArticlePersistsTombstone(t *testing.T) {
	home := t.TempDir()
	boardID := "test"
	boardDir := filepath.Join(home, "boards", "t", boardID)
	if err := os.MkdirAll(boardDir, 0750); err != nil {
		t.Fatal(err)
	}

	filename := "M.1723423500.A.123"
	record := pttbbs.NewFileHeader()
	record.SetFilename(filename)
	record.SetOwner("author")
	record.SetDate(" 8/12")
	record.SetTitle("delete me")
	record.AddRecommend(9)
	if err := pttbbs.AppendFileHeaderFileRecord(filepath.Join(boardDir, ".DIR"), record); err != nil {
		t.Fatal(err)
	}
	articlePath := filepath.Join(boardDir, filename)
	if err := os.WriteFile(articlePath, []byte("body"), 0600); err != nil {
		t.Fatal(err)
	}

	db, err := bbs.Open("pttbbs", home)
	if err != nil {
		t.Fatal(err)
	}
	repo := &repository{db: db}
	if err := repo.DeleteBoardArticle(context.Background(), boardID, filename, "author"); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(filepath.Join(boardDir, ".DIR")) // #nosec G304 -- path is under t.TempDir().
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) != 128 {
		t.Fatalf(".DIR size = %d, want 128 bytes", len(raw))
	}
	deleted, err := pttbbs.NewFileHeaderWithByte(raw)
	if err != nil {
		t.Fatal(err)
	}
	if deleted.Filename() != ".d" {
		t.Errorf("filename = %q, want .d", deleted.Filename())
	}
	if deleted.Owner() != "-" {
		t.Errorf("owner = %q, want -", deleted.Owner())
	}
	if deleted.Title() != "(本文已被刪除) [author]" {
		t.Errorf("title = %q", deleted.Title())
	}
	if deleted.Recommend() != 9 {
		t.Errorf("recommend = %d, want 9", deleted.Recommend())
	}
	if _, err := os.Stat(articlePath); !os.IsNotExist(err) {
		t.Fatalf("original article file still exists: %v", err)
	}
}
