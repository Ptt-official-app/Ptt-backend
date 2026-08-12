package repository

import (
	"testing"

	"github.com/Ptt-official-app/go-bbs/pttbbs"
)

func TestBoardReadPermissionLevel(t *testing.T) {
	board := &pttbbs.BoardHeader{Level: pttbbs.PermBM}

	level, ok := BoardReadPermissionLevel(board)
	if !ok {
		t.Fatal("expected PTT board permission level to be available")
	}
	if level != pttbbs.PermBM {
		t.Fatalf("expected level %d, got %d", pttbbs.PermBM, level)
	}

	board.Brdattr = pttbbs.BoardPostMask
	level, ok = BoardReadPermissionLevel(board)
	if !ok {
		t.Fatal("expected PTT post-mask board permission level to be available")
	}
	if level != 0 {
		t.Fatalf("post-mask board level is not a read restriction: got %d", level)
	}
}

func TestUserPermissionLevel(t *testing.T) {
	rawUser := &pttbbs.Userec{UserLevel: pttbbs.PermBM}

	level, ok := UserPermissionLevel(rawUser)
	if !ok {
		t.Fatal("expected PTT user permission level to be available")
	}
	if level != pttbbs.PermBM {
		t.Fatalf("expected level %d, got %d", pttbbs.PermBM, level)
	}

	wrappedUser := &bbsUserRecord{UserRecord: rawUser}
	level, ok = UserPermissionLevel(wrappedUser)
	if !ok {
		t.Fatal("expected wrapped PTT user permission level to be available")
	}
	if level != pttbbs.PermBM {
		t.Fatalf("expected wrapped level %d, got %d", pttbbs.PermBM, level)
	}
}
