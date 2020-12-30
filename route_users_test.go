package main

import (
	"encoding/json"
	"github.com/PichuChen/go-bbs"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetUserInformation(t *testing.T) {

	expected := bbs.Userec{
		Version:       4194,
		UserId:        "SYSOP",
		RealName:      "CodingMan",
		Nickname:      "神",
		Password:      "bhwvOJtfT1TAI",
		UserFlag:      0x02000A60,
		UserLevel:     0x20000407,
		NumLoginDays:  2,
		NumPosts:      0,
		FirstLogin:    time.Date(2020, 9, 21, 9, 41, 28, 0, time.UTC),
		LastLogin:     time.Date(2020, 9, 22, 6, 28, 14, 0, time.UTC),
		LastHost:      "59.124.167.226",
		Money:         0,
		Address:       "新竹縣子虛鄉烏有村543號",
		Over18:        true,
		Pager:         1,
		Invisible:     false,
		Career:        "全景軟體",
		LastSeen:      time.Date(2020, 9, 21, 9, 41, 28, 0, time.UTC),
		TimeSetAngel:  time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		TimePlayAngel: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		LastSong:      time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),

		TimeRemoveBadPost: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		TimeViolateLaw:    time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	userRecs = []*bbs.Userec{
		&expected,
	}

	req, err := http.NewRequest("GET", "/v1/users/SYSOP/information", nil)
	if err != nil {
		t.Fatal(err)
	}

	token := newAccessTokenWithUsername(expected.UserId)
	t.Logf("testing token: %v", token)
	req.Header.Add("Authorization", "bearer "+token)

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/users/", routeUsers)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responsedMap := map[string]interface{}{}
	json.Unmarshal(rr.Body.Bytes(), &responsedMap)
	t.Logf("got response %v", rr.Body.String())
	responsedData := responsedMap["data"].(map[string]interface{})
	if responsedData["user_id"] != expected.UserId {
		t.Errorf("handler returned unexpected body, user_id not match: got %v want userId %v",
			rr.Body.String(), expected.UserId)

	}

}

func TestParseUserPath(t *testing.T) {

	type TestCase struct {
		input         string
		expectdUserId string
		expectdItem   string
	}

	cases := []TestCase{
		{
			input:         "/v1/users/Pichu/information",
			expectdUserId: "Pichu",
			expectdItem:   "information",
		},
		{
			input:         "/v1/users/Pichu/",
			expectdUserId: "Pichu",
			expectdItem:   "",
		},
		{
			input:         "/v1/users/Pichu",
			expectdUserId: "Pichu",
			expectdItem:   "",
		},
	}

	for index, c := range cases {
		input := c.input
		expectdUserId := c.expectdUserId
		expectdItem := c.expectdItem
		actualUserId, actualItem, err := parseUserPath(input)
		if err != nil {
			t.Errorf("error on index %d, got: %v", index, err)

		}

		if actualUserId != expectdUserId {
			t.Errorf("userId not match on index %d, expected: %v, got: %v", index, expectdUserId, actualUserId)
		}

		if actualItem != expectdItem {
			t.Errorf("item not match on index %d, expected: %v, got: %v", index, expectdItem, actualItem)
		}

	}

}
