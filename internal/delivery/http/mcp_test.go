package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateMCPBoardID(t *testing.T) {
	for _, boardID := range []string{"Gossiping", "C_Chat", "NTU-Exam"} {
		if err := validateMCPBoardID(boardID); err != nil {
			t.Fatalf("validateMCPBoardID(%q): %v", boardID, err)
		}
	}
	for _, boardID := range []string{"", "../etc", "board/name", "board name"} {
		if err := validateMCPBoardID(boardID); err == nil {
			t.Fatalf("validateMCPBoardID(%q) unexpectedly succeeded", boardID)
		}
	}
}

func TestValidateMCPFilename(t *testing.T) {
	for _, filename := range []string{"M.1750603622.A.960", "12345"} {
		if err := validateMCPFilename(filename); err != nil {
			t.Fatalf("validateMCPFilename(%q): %v", filename, err)
		}
	}
	for _, filename := range []string{"", "../secret", "a/b", "file name"} {
		if err := validateMCPFilename(filename); err == nil {
			t.Fatalf("validateMCPFilename(%q) unexpectedly succeeded", filename)
		}
	}
}

func TestMCPHandlerRejectsLegacyGETStream(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/mcp", nil)

	(&Delivery{}).mcpHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET /mcp status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
}
