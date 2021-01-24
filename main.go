package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PichuChen/go-bbs"
	_ "github.com/PichuChen/go-bbs/pttbbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

var userRepo repository.UserRepository
var boardRepo repository.BoardRepository

var db *bbs.DB
var logger = logging.NewLogger()
var globalConfig *config.Config

func main() {
	logger.Informationalf("server start")

	var err error

	globalConfig, err = config.NewDefaultConfig()
	if err != nil {
		logger.Errorf("failed to get config: %v", err)
		return
	}

	db, err = bbs.Open("pttbbs", globalConfig.BBSHome)
	if err != nil {
		logger.Errorf("open bbs db error: %v", err)
		return
	}

	boardRepo, err = repository.NewBoardRepository(db)
	if err != nil {
		logger.Errorf("failed to create board repository: %s\n", err)
		return
	}

	userRepo, err = repository.NewUserRepository(db)
	if err != nil {
		logger.Errorf("failed to create user repository: %s\n", err)
		return
	}

	r := http.NewServeMux()
	buildRoute(r)

	logger.Informationalf("listen port on %v", globalConfig.ListenPort)
	err = http.ListenAndServe(fmt.Sprintf(":%v", globalConfig.ListenPort), r)
	if err != nil {
		logger.Errorf("listen serve error: %v", err)
	}
}

func routeClass(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getClass(w, r)
		return
	}

}

func getClass(w http.ResponseWriter, r *http.Request) {

	seg := strings.Split(r.URL.Path, "/")

	classId := "0"
	if len(seg) > 2 {
		classId = seg[3]
	}
	logger.Informationalf("user get class: %v", classId)

	list := []interface{}{}

	c := map[string]interface{}{
		"id":             1,
		"type":           "class",
		"title":          "title",
		"number_of_user": 3,
		"moderators": []string{
			"SYSOP",
			"pichu",
		},
	}
	list = append(list, c)

	m := map[string]interface{}{
		"data": list,
	}
	b, _ := json.MarshalIndent(m, "", "  ")

	w.Write(b)

}

// return a boolean value to indicate support guest account
// and using guest permission when permission insufficient
func supportGuest() bool {
	return false
}
