package main

import (
	"database/sql"
	"time"
)

type ArticleContent struct {
	BoardName string // board Name, e.g., "Gossiping"
	Filename  string // Filename, e.g., "M.1234567890.A.1BC"
	Content   []byte // Article content in bytes
	Extra     map[string]any
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateArticleContentTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS articles (
		board_name TEXT NOT NULL,
		filename TEXT NOT NULL,
		content BYTEA NOT NULL,
		extra JSONB DEFAULT '{}',
		created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
		updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
		PRIMARY KEY (board_name, filename)
	);
	`
	_, err := db.Exec(query)
	return err
}

func WriteArticleContentToDB(db *sql.DB, boardName string, filename string, content []byte) error {
	query := `
		INSERT INTO articles (board_name, filename, content)
		VALUES ($1, $2, $3)
		ON CONFLICT (board_name, filename) DO UPDATE
		SET content = EXCLUDED.content;
	`
	_, err := db.Exec(query, boardName, filename, content)
	return err
}

func ReadArticleContentFromDB(db *sql.DB, boardName string, filename string) ([]byte, error) {
	var content []byte
	query := `
		SELECT content
		FROM articles
		WHERE board_name = $1 AND filename = $2;
	`
	err := db.QueryRow(query, boardName, filename).Scan(&content)
	if err != nil {
		return nil, err
	}
	return content, nil
}
