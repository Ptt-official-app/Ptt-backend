package aids

import (
	"strings"
	"testing"
)

func TestFn2Aidu(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Aidu
	}{
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "invalid prefix",
			input:    "X.123456789.A.0DC",
			expected: 0,
		},
		{
			name:     "missing dot after prefix",
			input:    "M123456789.A.0DC",
			expected: 0,
		},
		{
			name:     "invalid timestamp",
			input:    "M.abc.A.0DC",
			expected: 0,
		},
		{
			name:     "timestamp is zero",
			input:    "M.0.A.0DC",
			expected: 0,
		},
		{
			name:     "missing .A after timestamp",
			input:    "M.123456789.B.0DC",
			expected: 0,
		},
		{
			name:     "missing random part",
			input:    "M.123456789.A.",
			expected: (0 << 44) | (123456789 << 12),
		},
		{
			name:     "invalid random part",
			input:    "M.123456789.A.GGG",
			expected: 0,
		},
		{
			name:     "random part is zero",
			input:    "M.123456789.A.000",
			expected: (0 << 44) | (123456789 << 12),
		},
		{
			name:     "valid M type",
			input:    "M.123456789.A.0DC",
			expected: (0 << 44) | (123456789 << 12) | 0x0DC,
		},
		{
			name:     "valid G type",
			input:    "G.987654321.A.1AB",
			expected: (1 << 44) | (987654321 << 12) | 0x1AB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Fn2Aidu(tt.input)
			if got != tt.expected {
				t.Errorf("Fn2Aidu(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}
func TestAidu2Aidc(t *testing.T) {
	tests := []struct {
		name     string
		input    Aidu
		expected string
	}{
		{
			name:     "zero value",
			input:    0,
			expected: "00000000",
		},
		{
			name:     "max value",
			input:    0xFFFFFFFFFFFF, // 48 bits set
			expected: "________",
		},
		{
			name:     "single digit",
			input:    1,
			expected: "00000001",
		},
		{
			name:     "two digits",
			input:    64,
			expected: "00000010",
		},
		{
			name:     "three digits",
			input:    65,
			expected: "00000011",
		},
		{
			name:     "multiple digits",
			input:    123456789,
			expected: "0007MyqL",
		},
		{
			name:     "random value",
			input:    (1 << 44) | (987654321 << 12) | 0x1AB,
			expected: Aidu2Aidc((1 << 44) | (987654321 << 12) | 0x1AB), // just check roundtrip
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Aidu2Aidc(tt.input)
			if len(got) != 8 {
				t.Errorf("Aidu2Aidc(%d) = %q, length = %d, want 8", tt.input, got, len(got))
			}
			// For the random value, check roundtrip
			if tt.name == "random value" {
				// No reverse function, so just check length and charset
				for i := 0; i < len(got); i++ {
					if !strings.ContainsRune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_", rune(got[i])) {
						t.Errorf("Aidu2Aidc(%d) = %q, contains invalid character %q", tt.input, got, got[i])
					}
				}
			} else if got != tt.expected {
				t.Errorf("Aidu2Aidc(%d) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
func TestAidu2Fn(t *testing.T) {
	tests := []struct {
		name     string
		input    Aidu
		expected string
	}{
		{
			name:     "zero value",
			input:    0,
			expected: "M.0.A.0",
		},
		{
			name:     "invalid type",
			input:    (2 << 44) | (123456789 << 12) | 0x0DC,
			expected: "",
		},
		{
			name:     "M type, random zero",
			input:    (0 << 44) | (123456789 << 12),
			expected: "M.123456789.A.0",
		},
		{
			name:     "M type, random nonzero",
			input:    (0 << 44) | (123456789 << 12) | 0x0DC,
			expected: "M.123456789.A.DC",
		},
		{
			name:     "G type, random nonzero",
			input:    (1 << 44) | (987654321 << 12) | 0x1AB,
			expected: "G.987654321.A.1AB",
		},
		{
			name:     "G type, random zero",
			input:    (1 << 44) | (987654321 << 12),
			expected: "G.987654321.A.0",
		},
		{
			name:     "max values",
			input:    (1 << 44) | (0xFFFFFFFF << 12) | 0xFFF,
			expected: "G.4294967295.A.FFF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Aidu2Fn(tt.input)
			if got != tt.expected {
				t.Errorf("Aidu2Fn(%d) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
func TestAidc2Aidu(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Aidu
	}{
		{
			name:     "all zeros",
			input:    "00000000",
			expected: 0,
		},
		{
			name:     "one zeros",
			input:    "0",
			expected: 0,
		},
		{
			name:     "all underscores (max value for 8 chars)",
			input:    "________",
			expected: 0xFFFFFFFFFFFF, // 48 bits set, max value for Aidc
		},
		{
			name:     "two underscores",
			input:    "__",
			expected: 63<<6 | 63, // 63 for underscore

		},
		{
			name:     "single digit",
			input:    "00000001",
			expected: 1,
		},
		{
			name:     "two digits",
			input:    "00000010",
			expected: 64,
		},
		{
			name:     "three digits",
			input:    "00000011",
			expected: 65,
		},
		{
			name:     "mixed case",
			input:    "0007MyqL",
			expected: 123456789,
		},
		{
			name:     "with dash",
			input:    "0000000-",
			expected: 62,
		},
		{
			name:     "with underscore",
			input:    "0000000_",
			expected: 63,
		},
		{
			name:     "invalid character",
			input:    "0000000@",
			expected: 0,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "random value roundtrip",
			input:    Aidu2Aidc((1 << 44) | (987654321 << 12) | 0x1AB),
			expected: (1 << 44) | (987654321 << 12) | 0x1AB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Aidc2Aidu(tt.input)
			if got != tt.expected {
				t.Errorf("Aidc2Aidu(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFn2Aidc(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal M type",
			input:    "M.1748774060.A.317",
			expected: "1eF2oiCN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Aidu2Aidc(Fn2Aidu(tt.input))
			if got != tt.expected {
				t.Errorf("Aidu2Aidc(Fn2Aidu(%q)) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestAidc2Fn(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal M type",
			input:    "1eF2oiCN",
			expected: "M.1748774060.A.317",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Aidu2Fn(Aidc2Aidu(tt.input))
			if got != tt.expected {
				t.Errorf("Aidu2Fn(Aidc2Aidu(%q)) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
