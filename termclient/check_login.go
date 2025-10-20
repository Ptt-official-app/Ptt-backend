package termclient

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

func isValidUsername(username string) bool {
	// Add your username validation logic here
	if len(username) == 0 {
		return false
	}
	for _, char := range username {
		if !(char >= 'a' && char <= 'z') &&
			!(char >= 'A' && char <= 'Z') &&
			!(char >= '0' && char <= '9') &&
			char != '_' {
			return false
		}
	}
	return true
}

func isValidPassword(password string) bool {
	// Add your password validation logic here
	if len(password) < 2 {
		return false
	}
	return true
}

func SetTermOutput(logger slog.Handler) {
	slog.SetDefault(slog.New(logger))
}

func (client *Client) CheckLogin(username, password string) (bool, error) {
	if !isValidUsername(username) {
		slog.Error("Invalid username format")
		return false, fmt.Errorf("invalid username format")
	}
	if !isValidPassword(password) {
		slog.Error("Invalid password format")
		return false, fmt.Errorf("invalid password format")
	}

	var 有重複登入 = false
	var 您有一篇文章尚未完成 = false
	var 請勿頻繁登入以免造成系統過度負荷 = false
	var 密碼不對或無此帳號請檢查大小寫及有無輸入錯誤 = false

	result := make(chan bool)

	go func() {
		// Script
		time.Sleep(time.Millisecond * 200)
		client.sendTextBySequence(username+"\r", time.Millisecond*10)
		client.sendTextBySequence(password+"\r", time.Millisecond*10)
		slog.Info("Password sent")
		<-time.After(time.Millisecond * 1200)
		if 有重複登入 || 您有一篇文章尚未完成 || 請勿頻繁登入以免造成系統過度負荷 {
			slog.Info("Duplicate login or unfinished article or frequent login warning detected")
			result <- true
		} else if 密碼不對或無此帳號請檢查大小寫及有無輸入錯誤 {
			slog.Info("Incorrect password or no such account detected")
			result <- false
		} else {
			slog.Info("No duplicate login detected")
		}
	}()
	go func() {
		for {
			t, msg, err := client.conn.ReadMessage()
			slog.Info("Read message", "type", t, "msg", len(msg), "err", err)
			if err != nil {
				slog.Info("Read message error", "err", err)
				return
			}

			// big5 to utf8
			big5Decoder := traditionalchinese.Big5.NewDecoder()
			utf8Data, _, err := transform.Bytes(big5Decoder, msg)
			if err != nil {
				utf8Data = msg
			}
			fmt.Printf("%s", string(utf8Data))
			if strings.Contains(string(utf8Data), "您有其它連線已登入此帳號") {
				有重複登入 = true
				slog.Info("Detected duplicate login")
			}
			if strings.Contains(string(utf8Data), "請勿頻繁登入以免造成系統過度負荷") {
				請勿頻繁登入以免造成系統過度負荷 = true
				slog.Info("Detected frequent login warning")
			}
			if strings.Contains(string(utf8Data), "上方為使用者心情點播留言區") {
				slog.Info("Detected on main menu")
			}
			if strings.Contains(string(utf8Data), "您有一篇文章尚未完成") {
				您有一篇文章尚未完成 = true
				slog.Info("Detected unfinished article")
			}
			if strings.Contains(string(utf8Data), "密碼不對或無此帳號。請檢查大小寫及有無輸入錯誤") {
				密碼不對或無此帳號請檢查大小寫及有無輸入錯誤 = true
				slog.Info("Detected incorrect password or no such account")
			}
		}
	}()
	r := <-result
	client.Disconnect()

	return r, nil
}
