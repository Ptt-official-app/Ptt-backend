package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
)

func TestLoginMetadataOverridesPersistentUserInformation(t *testing.T) {
	repo := &MockRepository{}
	uc := NewUsecase(&config.Config{}, repo)
	loginAt := time.Date(2026, time.August, 12, 3, 45, 0, 0, time.UTC)

	uc.recordLoginAt("pichu", "203.0.113.7", loginAt)

	data, err := uc.GetUserInformation(context.Background(), "pichu")
	if err != nil {
		t.Fatalf("GetUserInformation() error = %v", err)
	}
	if got := data["last_login_time"]; got != loginAt.Format(time.RFC3339) {
		t.Fatalf("last_login_time = %v, want %s", got, loginAt.Format(time.RFC3339))
	}
	if got := data["last_login_ip"]; got != "203.0.113.7" {
		t.Fatalf("last_login_ip = %v, want 203.0.113.7", got)
	}
	if got := data["last_login_ipv4"]; got != "203.0.113.7" {
		t.Fatalf("last_login_ipv4 = %v, want 203.0.113.7", got)
	}
}

func TestLoginMetadataUsesMostRecentRecord(t *testing.T) {
	repo := &MockRepository{}
	uc := NewUsecase(&config.Config{}, repo)

	uc.recordLoginAt("pichu", "192.0.2.1", time.Unix(10, 0))
	uc.recordLoginAt("pichu", "2001:db8::1", time.Unix(20, 0))

	data, err := uc.GetUserInformation(context.Background(), "pichu")
	if err != nil {
		t.Fatalf("GetUserInformation() error = %v", err)
	}
	if got := data["last_login_ip"]; got != "2001:db8::1" {
		t.Fatalf("last_login_ip = %v, want 2001:db8::1", got)
	}
}
