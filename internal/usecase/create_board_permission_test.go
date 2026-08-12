package usecase

import (
	"path/filepath"
	"testing"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/go-bbs"
	"github.com/Ptt-official-app/go-bbs/pttbbs"
)

func TestCreateBoardPermissionRequiresSYSOP(t *testing.T) {
	cfg, err := config.NewConfig("../../conf/config_default.toml", filepath.Join(t.TempDir(), "missing.toml"))
	if err != nil {
		t.Fatalf("load test config: %v", err)
	}

	for _, testCase := range []struct {
		name      string
		userLevel uint32
		wantErr   bool
	}{
		{name: "regular user denied", userLevel: 1, wantErr: true},
		{name: "SYSOP allowed", userLevel: pttbbs.PermSYSOP, wantErr: false},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			repo := &permissionTestRepository{
				MockRepository: &MockRepository{},
				users: []bbs.UserRecord{
					&permissionTestUser{MockUser: &MockUser{userID: "operator"}, level: testCase.userLevel},
				},
			}
			uc := NewUsecase(cfg, repo).(*usecase)
			token := uc.CreateAccessTokenWithUsername("operator")
			err := uc.CheckPermission(token, []Permission{PermissionCreateBoard}, nil)
			if testCase.wantErr && err == nil {
				t.Fatal("expected board creation permission error")
			}
			if !testCase.wantErr && err != nil {
				t.Fatalf("expected board creation permission, got %v", err)
			}
		})
	}
}
