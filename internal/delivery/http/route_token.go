package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Ptt-official-app/go-bbs"
)

// postToken will check username and verify password.
// It will generate the token return to client if it’s valid.
// Request URL: /v1/token
// API Doc: https://pttapp.cc/swagger/#/%E4%BD%BF%E7%94%A8%E8%80%85%E9%83%A8%E5%88%86/post_v1_token
func (delivery *Delivery) postToken(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err == nil {
		delivery.logger.Errorf("postToken parse form err: %w", err)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	userec, err := delivery.usecase.GetUserByID(context.Background(), username)
	if err != nil {
		m := map[string]string{
			"error":             "grant_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("postToken write get user id error response err: %w", err)
		}
		return

	}

	log.Println("found user:", userec)
	err = delivery.verifyPassword(userec, password)
	if err != nil {
		// TODO: add delay, warning, notify user

		m := map[string]string{
			"error":             "grant_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("postToken write password verify error response err: %w", err)
		}
		return
	}

	// Generate Access Token
	token := delivery.usecase.CreateAccessTokenWithUsername(username)
	m := map[string]string{
		"access_token": token,
		"token_type":   "bearer",
	}

	b, _ := json.MarshalIndent(m, "", "  ")

	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("postToken success response err: %w", err)
	}
}

func (delivery *Delivery) verifyPassword(userec bbs.UserRecord, password string) error {
	log.Println("password", userec.HashedPassword())
	return userec.VerifyPassword(password)
}

func (delivery *Delivery) getTokenFromRequest(r *http.Request) string {
	a := r.Header.Get("Authorization")
	s := strings.Split(a, " ")
	if len(s) < 2 {
		delivery.logger.Warningf("getTokenFromRequest error: len(s) < 2, got: %v", len(s))
		return ""
	}
	return s[1]
}
