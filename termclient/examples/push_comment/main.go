package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/termclient"
)

func main() {

	logFile, err := os.Create("log")
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()
	fileLogger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(fileLogger)

	username := os.Getenv("PTT_USERNAME")
	password := os.Getenv("PTT_PASSWORD")
	board := os.Getenv("PTT_BOARD")
	index := os.Getenv("PTT_INDEX")
	content_type := os.Getenv("PTT_COMMENT_CONTENT_TYPE")
	content := os.Getenv("PTT_COMMENT_CONTENT")

	// termclient.SetTermLogger(os.Stdout)
	client := termclient.NewClient()
	client.Connect()
	defer client.Disconnect()
	_, _ = client.PushComment(username, password, board, index, content_type, content)
	time.Sleep(3 * time.Second)
}
