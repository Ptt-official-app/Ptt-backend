package usecase

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/go-bbs"
)

type favoriteCaptureRepository struct {
	*MockRepository
	userID  string
	options bbs.FavoriteCreateOptions
}

func (repo *favoriteCaptureRepository) AddUserFavorite(_ context.Context, userID string, options bbs.FavoriteCreateOptions) (bbs.FavoriteRecord, error) {
	repo.userID = userID
	repo.options = options
	return nil, nil
}

func TestAddUserFavorite(t *testing.T) {
	tests := []struct {
		name         string
		tokenUserID  string
		targetUserID string
		favoriteType string
		boardID      string
		title        string
		wantType     bbs.FavoriteType
		wantErr      error
	}{
		{
			name:         "line",
			tokenUserID:  "pichu",
			targetUserID: "pichu",
			favoriteType: "line",
			wantType:     bbs.FavoriteTypeLine,
		},
		{
			name:         "folder",
			tokenUserID:  "pichu",
			targetUserID: "pichu",
			favoriteType: "folder",
			title:        "test",
			wantType:     bbs.FavoriteTypeFolder,
		},
		{
			name:         "board",
			tokenUserID:  "pichu",
			targetUserID: "pichu",
			favoriteType: "board",
			boardID:      "SYSOP",
			wantType:     bbs.FavoriteTypeBoard,
		},
		{
			name:         "cannot modify another user",
			tokenUserID:  "pichu",
			targetUserID: "someoneelse",
			favoriteType: "line",
			wantErr:      ErrFavoriteForbidden,
		},
		{
			name:         "folder title required",
			tokenUserID:  "pichu",
			targetUserID: "pichu",
			favoriteType: "folder",
			wantErr:      ErrInvalidFavorite,
		},
		{
			name:         "board id required",
			tokenUserID:  "pichu",
			targetUserID: "pichu",
			favoriteType: "board",
			wantErr:      ErrInvalidFavorite,
		},
		{
			name:         "invalid type",
			tokenUserID:  "pichu",
			targetUserID: "pichu",
			favoriteType: "unknown",
			wantErr:      ErrInvalidFavorite,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			globalConfig := favoriteTestConfig(t)
			repo := &favoriteCaptureRepository{MockRepository: &MockRepository{}}
			uc := NewUsecase(globalConfig, repo)
			token := uc.CreateAccessTokenWithUsername(tt.tokenUserID)

			_, err := uc.AddUserFavorite(context.Background(), token, tt.targetUserID, tt.favoriteType, tt.boardID, tt.title)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("AddUserFavorite() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("AddUserFavorite() error = %v", err)
			}
			if repo.userID != tt.targetUserID {
				t.Errorf("repository userID = %q, want %q", repo.userID, tt.targetUserID)
			}
			if repo.options.Type != tt.wantType {
				t.Errorf("favorite type = %v, want %v", repo.options.Type, tt.wantType)
			}
			if repo.options.BoardID != tt.boardID {
				t.Errorf("boardID = %q, want %q", repo.options.BoardID, tt.boardID)
			}
			if repo.options.Title != tt.title {
				t.Errorf("title = %q, want %q", repo.options.Title, tt.title)
			}
		})
	}
}

func TestAddUserFavoriteRejectsInvalidToken(t *testing.T) {
	repo := &favoriteCaptureRepository{MockRepository: &MockRepository{}}
	uc := NewUsecase(favoriteTestConfig(t), repo)

	_, err := uc.AddUserFavorite(context.Background(), "not-a-jwt", "pichu", "line", "", "")
	if !errors.Is(err, ErrFavoriteUnauthorized) {
		t.Fatalf("AddUserFavorite() error = %v, want unauthorized", err)
	}
}

func favoriteTestConfig(t *testing.T) *config.Config {
	t.Helper()
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	privateDER, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	publicDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	return &config.Config{
		AccessTokenPrivateKey: string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privateDER})),
		AccessTokenPublicKey:  string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER})),
		AccessTokenExpiresAt:  time.Hour,
	}
}
