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

func TestBoardIsPublic(t *testing.T) {
	testCases := []struct {
		name  string
		board *pttbbs.BoardHeader
		want  bool
	}{
		{
			name:  "normal board is public",
			board: &pttbbs.BoardHeader{BrdName: "Public"},
			want:  true,
		},
		{
			name:  "basic permission board is public",
			board: &pttbbs.BoardHeader{BrdName: "Basic", Level: 0o20},
			want:  true,
		},
		{
			name:  "hidden board is not public",
			board: &pttbbs.BoardHeader{BrdName: "Hidden", Brdattr: pttbbs.BoardHide},
			want:  false,
		},
		{
			name:  "top board is not public",
			board: &pttbbs.BoardHeader{BrdName: "Top", Brdattr: 0x00000800},
			want:  false,
		},
		{
			name:  "BM-only board is not public",
			board: &pttbbs.BoardHeader{BrdName: "BMOnly", Level: pttbbs.PermBM},
			want:  false,
		},
		{
			name: "post-mask level does not hide board",
			board: &pttbbs.BoardHeader{
				BrdName: "PostMask",
				Brdattr: pttbbs.BoardPostMask,
				Level:   pttbbs.PermBM,
			},
			want: true,
		},
		{
			name: "group board is not an article source",
			board: &pttbbs.BoardHeader{
				BrdName: "Group",
				Brdattr: pttbbs.BoardGroupBoard,
			},
			want: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := BoardIsPublic(testCase.board); got != testCase.want {
				t.Fatalf("BoardIsPublic() = %v, want %v", got, testCase.want)
			}
		})
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
