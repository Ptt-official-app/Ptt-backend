package main

import (
	"encoding/json"
	"fmt"
	"github.com/PichuChen/go-bbs"
	"net/http"
	"strings"
)

var userRecs []*bbs.Userec
var boardHeader []*bbs.BoardHeader

func main() {
	logger.Informationalf("server start")

	loadDefaultConfig()

	loadPasswdsFile()
	loadBoardFile()

	r := http.NewServeMux()
	buildRoute(r)

	logger.Informationalf("listen port on %v", globalConfig.ListenPort)
	err := http.ListenAndServe(fmt.Sprintf(":%v", globalConfig.ListenPort), r)
	if err != nil {
		logger.Errorf("listen serve error: %v", err)
	}
}

func loadPasswdsFile() {
	path, err := bbs.GetPasswdsPath(globalConfig.BBSHome)
	if err != nil {
		logger.Errorf("open file error: %v", err)
		return
	}
	logger.Debugf("path: %v", path)

	userRecs, err = bbs.OpenUserecFile(path)
	if err != nil {
		logger.Errorf("get user rec error: %v", err)
		return
	}
	logger.Debugf("userrec: %v", userRecs)
}

func loadBoardFile() {
	path, err := bbs.GetBoardPath(globalConfig.BBSHome)
	if err != nil {
		logger.Errorf("open file error: %v", err)
		return
	}
	logger.Debugf("path: %v", path)

	boardHeader, err = bbs.OpenBoardHeaderFile(path)
	if err != nil {
		logger.Errorf("get board header error: %v", err)
		return
	}
	// logger.Debugf("userrec: %v", userRecs)
	for index, board := range boardHeader {
		logger.Debugf("loaded %d %v", index, board.BrdName)

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
// and using guset permission when permission insufficient
func supportGuest() bool {
	return false
}
