package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	actual, err := NewConfig("testcase/01.toml", "")

	if err != nil {
		t.Errorf("NewConfig error excepted nil, got %v", err)
		return
	}

	expected := Config{
		BBSHome: "./home/bbs",
		// ListenPort            int16
		// AccessTokenPrivateKey string
		// AccessTokenPublicKey  string
		// AccessTokenExpiresAt  time.Duration
	}

	if actual.BBSHome != expected.BBSHome {
		t.Errorf("bbshome not match, expected: %v, got: %v", expected.BBSHome, actual.BBSHome)
		return
	}

}
