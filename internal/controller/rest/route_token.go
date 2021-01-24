package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PichuChen/go-bbs"
)

func (rest *restHandler) routeToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Check IP Flowspeed

	if r.Method == "POST" {
		rest.postToken(w, r)
		return
	}

}

func (rest *restHandler) postToken(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	username := r.FormValue("username")
	password := r.FormValue("password")

	userec, err := rest.findUserecById(username)
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
	err = rest.verifyPassword(userec, password)
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
	token := rest.getAccessTokenWithUsername(username)
	m := map[string]string{
		"access_token": token,
		"token_type":   "bearer",
	}

	b, _ := json.MarshalIndent(m, "", "  ")

	w.Write(b)

}

func (rest *restHandler) findUserecById(userid string) (bbs.UserRecord, error) {

	for _, it := range rest.userRepo.GetUsers(context.Background()) {
		if userid == it.UserId() {
			return it, nil
		}
	}
	return nil, fmt.Errorf("user record not found")

}

func (rest *restHandler) verifyPassword(userec bbs.UserRecord, password string) error {
	log.Println("password", userec.HashedPassword())
	return userec.VerifyPassword(password)
}

func (rest *restHandler) getTokenFromRequest(r *http.Request) string {
	a := r.Header.Get("Authorization")
	s := strings.Split(a, " ")
	if len(s) < 2 {
		rest.logger.Warningf("getTokenFromRequest error: len(s) < 2, got: %v", len(s))
		return ""
	}
	return s[1]
}
