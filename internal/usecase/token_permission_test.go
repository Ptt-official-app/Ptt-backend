package usecase

import (
	"context"
	"testing"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/go-bbs"
	"github.com/Ptt-official-app/go-bbs/pttbbs"
)

type permissionTestRepository struct {
	*MockRepository
	boards []bbs.BoardRecord
	users  []bbs.UserRecord
}

func (repo *permissionTestRepository) GetBoards(context.Context) []bbs.BoardRecord {
	return repo.boards
}

func (repo *permissionTestRepository) GetUsers(context.Context) ([]bbs.UserRecord, error) {
	return repo.users, nil
}

type permissionTestBoard struct {
	boardID string
	level   uint32
}

func (board *permissionTestBoard) BoardID() string         { return board.boardID }
func (board *permissionTestBoard) Title() string           { return board.boardID }
func (board *permissionTestBoard) IsClass() bool           { return false }
func (board *permissionTestBoard) ClassID() string         { return "" }
func (board *permissionTestBoard) BM() []string            { return nil }
func (board *permissionTestBoard) PermissionLevel() uint32 { return board.level }

type permissionTestUser struct {
	*MockUser
	level uint32
}

func (user *permissionTestUser) PermissionLevel() uint32 { return user.level }

func TestCheckBoardReadPermission(t *testing.T) {
	testCases := []struct {
		name          string
		requiredLevel uint32
		userLevel     uint32
		wantErr       bool
	}{
		{
			name:          "regular user cannot read BM-only board",
			requiredLevel: pttbbs.PermBM,
			userLevel:     1,
			wantErr:       true,
		},
		{
			name:          "BM can read BM-or-SYSOP board",
			requiredLevel: pttbbs.PermBM | pttbbs.PermSYSOP,
			userLevel:     pttbbs.PermBM,
			wantErr:       false,
		},
		{
			name:          "SYSOP can read BM-or-SYSOP board",
			requiredLevel: pttbbs.PermBM | pttbbs.PermSYSOP,
			userLevel:     pttbbs.PermSYSOP,
			wantErr:       false,
		},
		{
			name:          "board without required level remains readable",
			requiredLevel: 0,
			userLevel:     1,
			wantErr:       false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo := &permissionTestRepository{
				MockRepository: &MockRepository{},
				boards: []bbs.BoardRecord{
					&permissionTestBoard{boardID: "SECURITY", level: testCase.requiredLevel},
				},
				users: []bbs.UserRecord{
					&permissionTestUser{MockUser: &MockUser{userID: "pichu"}, level: testCase.userLevel},
				},
			}
			usecase := NewUsecase(&config.Config{}, repo).(*usecase)

			err := usecase.checkBoardReadPermission(context.Background(), "pichu", "SECURITY")
			if testCase.wantErr && err == nil {
				t.Fatal("expected permission error, got nil")
			}
			if !testCase.wantErr && err != nil {
				t.Fatalf("expected permission to be granted, got %v", err)
			}
		})
	}
}
