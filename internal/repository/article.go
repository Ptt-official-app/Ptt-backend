package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/go-bbs"
)

// PopularArticleRecord is an ArticleRecord which has boardID information.
type PopularArticleRecord interface {
	// Note: go-bbs has not implemented this yet
	// TODO: use bbs.PopularArticleRecord or something when it is ready
	bbs.ArticleRecord
	BoardID() string
}

type PopularArticle struct {
	filename       string
	modified       time.Time
	recommendCount int
	owner          string
	date           string
	title          string
	money          int
	boardID        string
}

func (p *PopularArticle) Filename() string               { return p.filename }
func (p *PopularArticle) Modified() time.Time            { return p.modified }
func (p *PopularArticle) SetModified(newValue time.Time) { p.modified = newValue }
func (p *PopularArticle) Recommend() int                 { return p.recommendCount }
func (p *PopularArticle) Date() string                   { return p.date }
func (p *PopularArticle) Title() string                  { return p.title }
func (p *PopularArticle) Money() int                     { return p.money }
func (p *PopularArticle) Owner() string                  { return p.owner }
func (p *PopularArticle) BoardID() string                { return p.boardID }

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
	var result []PopularArticleRecord
	boards := repo.GetBoards(ctx)
	for _, board := range boards {
		articles, err := repo.GetBoardArticleRecords(ctx, board.BoardID())
		if err != nil {
			return nil, err
		}
		for _, article := range articles {
			result = append(result, &PopularArticle{
				article.Filename(),
				article.Modified(),
				article.Recommend(),
				article.Owner(),
				article.Date(),
				article.Title(),
				article.Money(),
				board.BoardID(),
			})
		}
	}
	return result, nil
}

func (repo *repository) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (PushRecord, error) {
	// Append comment into board article file
	now := time.Now()
	result := "[1;31m " + appendType + " " + userID + " [m [33m: " + text + " [m " + now.Format("01/02 15:04") + "\n"
	err := repo.db.AppendBoardArticleFile(boardID, filename, bbs.Utf8ToBig5(result))
	if err != nil {
		return nil, err
	}

	p := &Push{
		appendType: appendType,
		id:         userID,
		ipAddr:     "", // not sure how to get IPAddr
		text:       result,
		time:       now,
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
	currentTime := time.Now()
	now := currentTime.Format("01/02")
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

	t, err := template.New("Ptt-article-template").Parse(config.PttArticleTemplate)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	err = t.Execute(buffer, config.ArticleArguments{
		UserData:      userData,
		Article:       record,
		BoardID:       boardID,
		Content:       content,
		PostANSICDate: currentTime.Format(time.ANSIC),
	})
	if err != nil {
		return nil, err
	}

	err = repo.db.WriteBoardArticleFile(boardID, record.Filename(), bbs.Utf8ToBig5(buffer.String()))
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
