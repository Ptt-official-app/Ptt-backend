package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Ptt-official-app/go-bbs"
)

// PopularArticleRecord is an ArticleRecord which has boardID information.
type PopularArticleRecord interface {
	// Note: go-bbs has not implemented this yet
	// TODO: use bbs.PopularArticleRecord or something when it is ready
	bbs.ArticleRecord
	BoardID() string
}

type PushRecord interface {
	// TODO: use bbs.PushRecord instead
	Type() string
	ID() string
	IPAddr() string
	Text() string
	Time() time.Time
}

func (repo *repository) GetPopularArticles(ctx context.Context) ([]PopularArticleRecord, error) {
	// Note: go-bbs has not implemented this yet
	// TODO: delegate to repo.db when it is ready
	return []PopularArticleRecord{}, nil
}

func (repo *repository) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (PushRecord, error) {
	return nil, nil
}

func (repo *repository) AppendArticle(ctx context.Context, userID, boardID, title, content string) (bbs.ArticleRecord, error) {
	return nil, nil
}

// CreateArticle
// TODO: return result from bbs response
func (repo *repository) CreateArticle(ctx context.Context, userID, boardID, title, content string) (bbs.ArticleRecord, error) {
	// get file name
	now := time.Now().Format("01/02")
	record, err := repo.db.CreateArticleRecord(boardID, userID, now, title)
	if err != nil {
		fmt.Println("CreateArticleRecord error:", err)
		return nil, err
	}

	err = repo.db.AddArticleRecordFileRecord(boardID, record)
	if err != nil {
		fmt.Println("AddArticleRecordFileRecord error:", err)
		return nil, err

	}

	err = repo.db.WriteBoardArticleFile(boardID, record.Filename(), bbs.Utf8ToBig5(content))
	if err != nil {
		fmt.Println("WriteBoardArticleFile error: %w", err)
		return nil, err
	}

	return record, nil
}

func (repo *repository) GetRawArticle(boardID, filename string) (string, error) {
	data, err := repo.db.ReadBoardArticleFile(boardID, filename)

	if err != nil {
		fmt.Println("ReadrBoardArticleFile error: %w", err)
		return "", err
	}

	return bbs.Big5ToUtf8(data), nil
}

type ForwardArticleToBoardRecord interface {
	// Note: go-bbs has not implemented this yet
	// TODO: use bbs.ForwardArticleToBoardRecord or something when it is ready
	bbs.ArticleRecord
	DestBoardID() string
	IPAddr() string
	ForwardTime() time.Time
	ForwardTitle() string
}

func (repo *repository) ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (ForwardArticleToBoardRecord, error) {
	// Note: go-bbs has not implemented this yet
	// TODO: delegate to repo.db when it is ready
	return nil, nil
}
