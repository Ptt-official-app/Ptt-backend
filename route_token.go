package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PichuChen/go-bbs"
)

func routeToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed
	if r.Method == http.MethodPost {
		postToken(w, r)
		return
	}
}

func postToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("failed to parse form data")
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	userec, err := findUserecByID(username)
	if err != nil {
		m := map[string]string{
			"error":             "grant_error",
			"error_description": err.Error(),
		}
		b, err := json.MarshalIndent(m, "", "  ")

		if err != nil {
			logger.Errorf("failed to marshal response data: %s\n", err)
		}

		if _, err := w.Write(b); err != nil {
			logger.Errorf("failed to write response: %s\n", err)
		}

		return
	}

	log.Println("found user:", userec)
	err = verifyPassword(userec, password)

	if err != nil {
		// TODO: add delay, warning, notify user
		m := map[string]string{
			"error":             "grant_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")

		w.WriteHeader(http.StatusUnauthorized)

		if _, err := w.Write(b); err != nil {
			logger.Errorf("failed to write response: %s\n", err)
		}

		return
	}

	// Generate Access Token
	token := newAccessTokenWithUsername(username)
	m := map[string]string{
		"access_token": token,
		"token_type":   "bearer",
	}

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		logger.Errorf("failed to marshal response data: %s\n", err)
	}

	if _, err := w.Write(b); err != nil {
		logger.Errorf("failed to write response: %s\n", err)
	}
}

func findUserecByID(userid string) (bbs.UserRecord, error) {
	for _, it := range userRecs {
		if userid == it.UserId() {
			return it, nil
		}
	}

	return nil, fmt.Errorf("user record not found")
}

func verifyPassword(userec bbs.UserRecord, password string) error {
	log.Println("password", userec.HashedPassword())
	return userec.VerifyPassword(password)
}

func getTokenFromRequest(r *http.Request) string {
	a := r.Header.Get("Authorization")
	s := strings.Split(a, " ")

	if len(s) < 2 {
		logger.Warningf("getTokenFromRequest error: len(s) < 2, got: %v", len(s))
		return ""
	}

	return s[1]
}
