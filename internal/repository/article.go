package repository

import (
	"context"
	"errors"
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

type Push struct {
	appendType string
	id         string
	ipAddr     string
	text       string
	time       time.Time
}

func (p *Push) Type() string {
	return p.appendType
}

func (p *Push) ID() string {
	return p.id
}

func (p *Push) IPAddr() string {
	return p.ipAddr
}

func (p *Push) Text() string {
	return p.text
}

func (p *Push) Time() time.Time {
	return p.time
}

func (repo *repository) GetPopularArticles(ctx context.Context) ([]PopularArticleRecord, error) {
	// Note: go-bbs has not implemented this yet
	// TODO: delegate to repo.db when it is ready
	return []PopularArticleRecord{}, nil
}

func (repo *repository) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (PushRecord, error) {
	// Append comment into board article file
	err := repo.db.AppendBoardArticleFile(filename, bbs.Utf8ToBig5(text))
	if err != nil {
		return nil, err
	}

	p := &Push{
		appendType: appendType,
		id:         userID,
		ipAddr:     "", // not sure how to get IPAddr
		text:       text,
		time:       time.Now(),
	}
	return p, nil
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

	var userData bbs.UserRecord = nil
	records, err := repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		if record.UserID() == userID {
			userData = record
		}
	}
	if userData == nil {
		return nil, errors.New("user ID not found")
	}

	err = repo.db.WriteBoardArticleFile(boardID, record.Filename(), bbs.Utf8ToBig5(fmt.Sprintf(`作者: %s (%s) 看板: %s
標題: %s
時間: %s


%s

--
※ 發信站: 新批踢踢(ptt2.cc), 來自: %s
※ 文章網址: http://www.ptt.cc/bbs/%s/%s.html
`, userID, userData.Nickname(), boardID,
		title,
		time.Now().Format(time.ANSIC),
		content,
		userData.LastHost(),
		boardID, record.Filename(),
	)))
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
