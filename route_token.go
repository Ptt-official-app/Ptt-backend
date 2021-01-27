package main

import (
	"encoding/json"
	"fmt"

	"github.com/PichuChen/go-bbs"
	"github.com/julienschmidt/httprouter"

	"log"
	"net/http"
	"strings"
)

func postToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: Check IP Flowspeed
	r.ParseForm()

	username := r.FormValue("username")
	password := r.FormValue("password")

	userec, err := findUserecById(username)
	if err != nil {
		m := map[string]string{
			"error":             "grant_error",
			"error_description": err.Error(),
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		w.Write(b)
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
		w.Write(b)
		return
	}

	// Generate Access Token
	token := newAccessTokenWithUsername(username)
	m := map[string]string{
		"access_token": token,
		"token_type":   "bearer",
	}

	b, _ := json.MarshalIndent(m, "", "  ")

	w.Write(b)

}

func findUserecById(userid string) (bbs.UserRecord, error) {

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
