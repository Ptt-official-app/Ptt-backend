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
)

var userRecs []bbs.UserRecord
var boardHeader []bbs.BoardRecord

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

	loadPasswdsFile()
	loadBoardFile()

	r := http.NewServeMux()
	buildRoute(r)

	logger.Informationalf("listen port on %v", globalConfig.ListenPort)
	err = http.ListenAndServe(fmt.Sprintf(":%v", globalConfig.ListenPort), r)
	if err != nil {
		logger.Errorf("listen serve error: %v", err)
	}
}

func loadPasswdsFile() {
	var err error
	userRecs, err = db.ReadUserRecords()
	if err != nil {
		logger.Errorf("get user rec error: %v", err)
		return
	}
}

func loadBoardFile() {
	var err error
	boardHeader, err = db.ReadBoardRecords()
	if err != nil {
		logger.Errorf("get board header error: %v", err)
		return
	}
	for index, board := range boardHeader {
		logger.Debugf("loaded %d %v", index, board.BoardId())

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
