package main

import (
	"encoding/json"
	"fmt"
	"github.com/PichuChen/go-bbs"
	"github.com/PichuChen/go-bbs/crypt"

	"log"
	"net/http"
	"strings"
)

func routeToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed

	if r.Method == "POST" {
		postToken(w, r)
		return
	}

}

func postToken(w http.ResponseWriter, r *http.Request) {

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

	log.Println("found user:", *userec)
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

func findUserecById(userid string) (*bbs.Userec, error) {

	for _, it := range userRecs {
		if userid == it.UserId {
			return it, nil
		}
	}
	return nil, fmt.Errorf("user record not found")

}

func verifyPassword(userec *bbs.Userec, password string) error {
	res, err := crypt.Fcrypt([]byte(password), []byte(userec.Password[:2]))
	if err != nil {
		logger.Criticalf("crypt.Fcrypt error: %v", err)
		return err

	}
	str := strings.Trim(string(res), "\x00")

	if str != userec.Password {
		return fmt.Errorf("password incorrect")
	}
	return nil

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
