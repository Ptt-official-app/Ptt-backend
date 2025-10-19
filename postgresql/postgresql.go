package postgresql

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	_ "github.com/lib/pq"

	"github.com/Ptt-official-app/go-bbs"
)

type Connector struct {
	DSN string
	db  *sql.DB
}

func init() {
	// register driver
	bbs.Register("postgresql", &Connector{})
}

// Open connect a postgresql database
// for example: postgresql://user:password@localhost:5432/database?sslmode=enable
func (c *Connector) Open(dataSourceName string) error {
	c.DSN = dataSourceName
	// parse DSN
	if !strings.HasPrefix(dataSourceName, "postgresql://") {
		return fmt.Errorf("invalid postgresql DSN: %s, "+
			"should be postgresql://user:password@localhost:5432/database?sslmode=enable", dataSourceName)
	}

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return fmt.Errorf("open postgresql error: %w", err)
	}
	c.db = db
	return nil
}
func (c *Connector) GetUserRecordsPath() (string, error) {
	return "user_record", nil
}

func (c *Connector) ReadUserRecordsFile(filename string) ([]bbs.UserRecord, error) {
	slog.Info("ReadUserRecordsFile", "filename", filename)
	return []bbs.UserRecord{}, nil
}
func (c *Connector) GetUserDraftPath(userID, draftID string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c *Connector) GetUserFavoriteRecordsPath(userID string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c *Connector) ReadUserFavoriteRecordsFile(filename string) ([]bbs.FavoriteRecord, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Connector) GetBoardRecordsPath() (string, error) {
	return "boards_table", nil
}

func (c *Connector) ReadBoardRecordsFile(path string) ([]bbs.BoardRecord, error) {
	return []bbs.BoardRecord{}, nil
}

func (c *Connector) GetBoardArticleRecordsPath(boardID string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c *Connector) ReadArticleRecordsFile(filename string) ([]bbs.ArticleRecord, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Connector) GetBoardTreasureRecordsPath(boardID string, treasureID []string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c *Connector) GetBoardArticleFilePath(boardID string, filename string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c *Connector) GetBoardTreasureFilePath(boardID string, treasureID []string, filename string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c *Connector) ReadBoardArticleFile(filename string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}
