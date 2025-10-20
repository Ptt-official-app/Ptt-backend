package termclient

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

// PushComment logs in (if needed), goes to board, opens the post by index, and sends a comment.
// contentType accepts: "推"/"1", "噓"/"2", "→"/"3" (case-insensitive for ascii fallbacks like PUSH/BOO/ARROW).
// Returns true if the flow completed without obvious error prompts, false otherwise.
func (client *Client) PushComment(username, password, board, index, contentType, content string) (bool, error) {
	if client.conn == nil {
		return false, fmt.Errorf("not connected")
	}
	if !isValidUsername(username) {
		slog.Error("Invalid username format")
		return false, fmt.Errorf("invalid username format")
	}
	if !isValidPassword(password) {
		slog.Error("Invalid password format")
		return false, fmt.Errorf("invalid password format")
	}
	if strings.TrimSpace(board) == "" || strings.TrimSpace(index) == "" {
		return false, fmt.Errorf("board and index are required")
	}

	// Normalize contentType to option number string: 1=推, 2=噓, 3=→
	pick := func(ct string) string {
		c := strings.TrimSpace(ct)
		switch c {
		case "1", "推":
			return "1"
		case "2", "噓":
			return "2"
		case "3", "→":
			return "3"
		}
		lc := strings.ToLower(c)
		if lc == "push" || lc == "p" {
			return "1"
		}
		if lc == "boo" || lc == "b" || lc == "boo!" {
			return "2"
		}
		if lc == "arrow" || lc == "a" || lc == "->" || lc == ">" {
			return "3"
		}
		return "1" // default to 推
	}
	opt := pick(contentType)

	// Flags updated by reader goroutine
	var 有重複登入 bool
	var 您有一篇文章尚未完成 bool
	var 請勿頻繁登入以免造成系統過度負荷 bool
	var 密碼不對或無此帳號請檢查大小寫及有無輸入錯誤 bool
	var 作者本人 bool

	// additional state flags that affect flow handled inline in reader

	done := make(chan bool, 1)
	failed := make(chan error, 1)

	// Reader: scan output and set state flags
	go func() {
		for {
			t, msg, err := client.conn.ReadMessage()
			slog.Info("Read message", "type", t, "len", len(msg), "err", err)
			if err != nil {
				failed <- err
				return
			}
			// big5 to utf8
			big5Decoder := traditionalchinese.Big5.NewDecoder()
			utf8Data, _, err := transform.Bytes(big5Decoder, msg)
			if err != nil {
				utf8Data = msg
			}
			m := string(utf8Data)
			fmt.Print(m)

			// Login related
			if strings.Contains(m, "您有其它連線已登入此帳號") {
				有重複登入 = true
				slog.Info("Detected duplicate login")
			}
			if strings.Contains(m, "請勿頻繁登入以免造成系統過度負荷") {
				請勿頻繁登入以免造成系統過度負荷 = true
				slog.Info("Detected frequent login warning")
			}
			if strings.Contains(m, "您有一篇文章尚未完成") {
				您有一篇文章尚未完成 = true
				slog.Info("Detected unfinished article")
			}
			if strings.Contains(m, "密碼不對或無此帳號。請檢查大小寫及有無輸入錯誤") {
				密碼不對或無此帳號請檢查大小寫及有無輸入錯誤 = true
				slog.Info("Detected incorrect password or no such account")
			}
			// Comment prompts / errors
			if strings.Contains(m, "禁止快速連續推文") {
				failed <- fmt.Errorf("禁止快速連續推文")
				return
			}
			if strings.Contains(m, "禁止短時間內大量推文") {
				failed <- fmt.Errorf("禁止短時間內大量推文")
				return
			}
			if strings.Contains(m, "使用者不可發言") {
				failed <- fmt.Errorf("使用者不可發言")
				return
			}
			if strings.Contains(m, "◆ 抱歉, 禁止推薦") {
				failed <- fmt.Errorf("抱歉, 禁止推薦")
				return
			}
			if strings.Contains(m, "作者本人, 使用 → 加註方式") {
				作者本人 = true
				slog.Info("Detected author self-commenting prompt")
			}
		}
	}()

	// Script: login → board → open article → comment
	go func() {
		// Login
		time.Sleep(200 * time.Millisecond)
		client.sendTextBySequence(username+"\r", 10*time.Millisecond)
		client.sendTextBySequence(password+"\r", 10*time.Millisecond)
		slog.Info("Credentials sent")

		// Allow server screens to settle
		time.Sleep(1200 * time.Millisecond)
		if 密碼不對或無此帳號請檢查大小寫及有無輸入錯誤 {
			done <- false
			return
		}

		// Handle duplicate login prompt (best-effort)
		if 有重複登入 {
			client.sendText("n")
			time.Sleep(200 * time.Millisecond)
			client.sendText("\r\n")
			time.Sleep(800 * time.Millisecond)
		}

		// Handle frequent login warning
		if 請勿頻繁登入以免造成系統過度負荷 {
			client.sendText("\r\n")
			time.Sleep(1500 * time.Millisecond)
		}

		// Continue screens
		client.sendText("\r\n")
		time.Sleep(1500 * time.Millisecond)

		if 您有一篇文章尚未完成 {
			// Quit unfinished draft
			client.sendText("Q\r\n")
			time.Sleep(1200 * time.Millisecond)
		}

		// Navigate to board
		slog.Info("Navigating to board", "board", board)
		client.sendText("s")
		time.Sleep(800 * time.Millisecond)
		client.sendTextBySequence(board+"\r", 10*time.Millisecond)
		time.Sleep(1200 * time.Millisecond)

		// Try to enter article by index
		client.sendTextBySequence(index+"\r", 10*time.Millisecond)
		time.Sleep(1000 * time.Millisecond)

		// Start comment
		client.sendText("X")
		time.Sleep(800 * time.Millisecond)

		// Choose comment type if asked
		if !作者本人 {
			client.sendTextBySequence(opt, 10*time.Millisecond)
			client.sendText("\r")
			time.Sleep(400 * time.Millisecond)
		}

		// Input content
		if strings.TrimSpace(content) != "" {
			client.sendTextBySequence(content, 10*time.Millisecond)
		}
		client.sendText("\r")
		time.Sleep(400 * time.Millisecond)

		// Confirm send
		client.sendText("y")
		client.sendText("\r")
		time.Sleep(1200 * time.Millisecond)

		// If no error flags raised by now, assume success
		done <- true
	}()

	// Wait for completion or failure with timeout
	select {
	case ok := <-done:
		if !ok {
			client.Disconnect()
			return false, fmt.Errorf("login failed")
		}
		return true, nil
	case err := <-failed:
		client.Disconnect()
		return false, err
	case <-time.After(20 * time.Second):
		client.Disconnect()
		return false, fmt.Errorf("push comment timeout")
	}
}
