package usecase

import (
	"testing"

	"github.com/Ptt-official-app/go-bbs/pttbbs"
)

type deletePermissionUser struct {
	*MockUser
	level uint32
}

func (u *deletePermissionUser) PermissionLevel() uint32 { return u.level }

type deletePermissionBoard struct {
	*MockBoardRecord
	bms    []string
	noSelf bool
}

func (b *deletePermissionBoard) BM() []string             { return b.bms }
func (b *deletePermissionBoard) IsHide() bool             { return false }
func (b *deletePermissionBoard) IsPostMask() bool         { return false }
func (b *deletePermissionBoard) IsAnonymous() bool        { return false }
func (b *deletePermissionBoard) IsDefaultAnonymous() bool { return false }
func (b *deletePermissionBoard) IsNoCredit() bool         { return false }
func (b *deletePermissionBoard) IsVoteBoard() bool        { return false }
func (b *deletePermissionBoard) IsWarnEL() bool           { return false }
func (b *deletePermissionBoard) IsTop() bool              { return false }
func (b *deletePermissionBoard) IsNoRecommend() bool      { return false }
func (b *deletePermissionBoard) IsAngelAnonymous() bool   { return false }
func (b *deletePermissionBoard) IsBMCount() bool          { return false }
func (b *deletePermissionBoard) IsNoBoo() bool            { return false }
func (b *deletePermissionBoard) IsRestrictedPost() bool   { return false }
func (b *deletePermissionBoard) IsGuestPost() bool        { return false }
func (b *deletePermissionBoard) IsCooldown() bool         { return false }
func (b *deletePermissionBoard) IsCPLog() bool            { return false }
func (b *deletePermissionBoard) IsNoFastRecommend() bool  { return false }
func (b *deletePermissionBoard) IsIPLogRecommend() bool   { return false }
func (b *deletePermissionBoard) IsOver18() bool           { return false }
func (b *deletePermissionBoard) IsNoReply() bool          { return false }
func (b *deletePermissionBoard) IsAlignedComment() bool   { return false }
func (b *deletePermissionBoard) IsNoSelfDeletePost() bool { return b.noSelf }
func (b *deletePermissionBoard) IsBMMaskContent() bool    { return false }

func TestCanDeleteArticle(t *testing.T) {
	article := &MockArticleRecord{filename: "M.1.A.001", owner: "author"}

	tests := []struct {
		name   string
		user   *deletePermissionUser
		board  *deletePermissionBoard
		wantOK bool
	}{
		{
			name:   "sysop can delete another user's article",
			user:   &deletePermissionUser{MockUser: &MockUser{userID: "SYSOP"}, level: pttbbs.PermSYSOP},
			board:  &deletePermissionBoard{MockBoardRecord: NewMockBoardRecord("1", "test", "test", false)},
			wantOK: true,
		},
		{
			name:   "board moderator can delete another user's article",
			user:   &deletePermissionUser{MockUser: &MockUser{userID: "moderator"}},
			board:  &deletePermissionBoard{MockBoardRecord: NewMockBoardRecord("1", "test", "test", false), bms: []string{"moderator"}},
			wantOK: true,
		},
		{
			name:   "author can delete own article",
			user:   &deletePermissionUser{MockUser: &MockUser{userID: "author"}},
			board:  &deletePermissionBoard{MockBoardRecord: NewMockBoardRecord("1", "test", "test", false)},
			wantOK: true,
		},
		{
			name:   "board can disable author self deletion",
			user:   &deletePermissionUser{MockUser: &MockUser{userID: "author"}},
			board:  &deletePermissionBoard{MockBoardRecord: NewMockBoardRecord("1", "test", "test", false), noSelf: true},
			wantOK: false,
		},
		{
			name:   "unrelated user cannot delete article",
			user:   &deletePermissionUser{MockUser: &MockUser{userID: "stranger"}},
			board:  &deletePermissionBoard{MockBoardRecord: NewMockBoardRecord("1", "test", "test", false)},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canDeleteArticle(tt.user, tt.board, article); got != tt.wantOK {
				t.Fatalf("canDeleteArticle() = %v, want %v", got, tt.wantOK)
			}
		})
	}
}
