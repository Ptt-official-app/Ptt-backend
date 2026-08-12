package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Ptt-official-app/go-bbs"
	"github.com/Ptt-official-app/go-bbs/pttbbs"
)

func TestAddUserFavoritePersistsAndReadsBack(t *testing.T) {
	home := t.TempDir()
	board := pttbbs.NewBoardHeader()
	board.SetBoardID("SYSOP")
	board.SetTitle("system board")
	boardPath, err := pttbbs.GetBoardPath(home)
	if err != nil {
		t.Fatal(err)
	}
	if err := pttbbs.AppendBoardHeaderFileRecord(boardPath, board); err != nil {
		t.Fatal(err)
	}

	favoritePath, err := pttbbs.GetUserFavoritePath(home, "pichu")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(favoritePath), 0750); err != nil {
		t.Fatal(err)
	}

	db, err := bbs.Open("pttbbs", home)
	if err != nil {
		t.Fatal(err)
	}
	repo := &repository{db: db}
	ctx := context.Background()

	if _, err := repo.AddUserFavorite(ctx, "pichu", bbs.FavoriteCreateOptions{Type: bbs.FavoriteTypeLine}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.AddUserFavorite(ctx, "pichu", bbs.FavoriteCreateOptions{Type: bbs.FavoriteTypeFolder, Title: "test"}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.AddUserFavorite(ctx, "pichu", bbs.FavoriteCreateOptions{Type: bbs.FavoriteTypeBoard, BoardID: "SYSOP"}); err != nil {
		t.Fatal(err)
	}

	records, err := repo.GetUserFavoriteRecords(ctx, "pichu")
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 3 {
		t.Fatalf("favorite count = %d, want 3", len(records))
	}
	if records[0].Type() != bbs.FavoriteTypeLine {
		t.Errorf("favorite[0] type = %v, want line", records[0].Type())
	}
	if records[1].Type() != bbs.FavoriteTypeFolder || records[1].Title() != "test" {
		t.Errorf("favorite[1] = type %v title %q", records[1].Type(), records[1].Title())
	}
	if records[2].Type() != bbs.FavoriteTypeBoard || records[2].BoardID() != "SYSOP" {
		t.Errorf("favorite[2] = type %v board %q", records[2].Type(), records[2].BoardID())
	}

	if info, err := os.Stat(favoritePath); err != nil {
		t.Fatalf("favorite file not persisted: %v", err)
	} else if info.Size() == 0 {
		t.Fatal("favorite file is empty")
	}

	// ReadUserFavoriteRecords resolves board IDs through .BRD. On Windows this
	// remove fails if OpenBoardHeaderFile leaked its handle.
	if err := os.Remove(boardPath); err != nil {
		t.Fatalf("board file remains open after favorite read: %v", err)
	}
}
