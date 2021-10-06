package http

import (
	"context"
	"time"

	"github.com/Ptt-official-app/go-bbs"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

func (usecase *MockUsecase) GetPopularArticles(ctx context.Context) ([]repository.PopularArticleRecord, error) {
	return []repository.PopularArticleRecord{}, nil
}

func (usecase *MockUsecase) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (repository.PushRecord, error) {
	return &MockPushRecord{
		appendType : "推",
		id : "test",
		ipAddr : "127.0.0.1",
		text : "test push",
		time : time.Now(),
	}, nil
}

func (usecase *MockUsecase) ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (repository.ForwardArticleToBoardRecord, error) {
	return nil, nil
}

func (usecase *MockUsecase) ForwardArticleToEmail(ctx context.Context, userID, boardID, filename, email string) error {
	return nil
}

func (usecase *MockUsecase) CreateArticle(ctx context.Context, userID, boardID, title, article string) (bbs.ArticleRecord, error) {
	return &MockArticleRecord{
		filename:       "filename1",
		modified:       time.Time{},
		recommendCount: 10,
		owner:          "SYSOP",
		date:           "",
		title:          "[討論] 偶爾要發個廢文",
		money:          0,
	}, nil
}

func (usecase *MockUsecase) GetRawArticle(boardID, filename string) (string, error) {
	return "test", nil
}

func (usecase *MockUsecase) UpdateUsefulness(ctx context.Context, userID, boardID, filename, appendType string) (repository.PushRecord, error) {
	return &MockPushRecord{
		appendType : "推",
		id : "test",
		ipAddr : "127.0.0.1",
		text : "test push",
		time : time.Now(),
	}, nil
}

type MockArticleRecord struct {
	filename       string
	modified       time.Time
	recommendCount int
	owner          string
	date           string
	title          string
	money          int
}

func (m *MockArticleRecord) Filename() string               { return m.filename }
func (m *MockArticleRecord) Modified() time.Time            { return m.modified }
func (m *MockArticleRecord) SetModified(newValue time.Time) { m.modified = newValue }
func (m *MockArticleRecord) Recommend() int                 { return m.recommendCount }
func (m *MockArticleRecord) Owner() string                  { return m.owner }
func (m *MockArticleRecord) Date() string                   { return m.date }
func (m *MockArticleRecord) Title() string                  { return m.title }
func (m *MockArticleRecord) Money() int                     { return m.money }

type MockPushRecord struct {
	appendType string
	id         string
	ipAddr     string
	text       string
	time       time.Time
}

func (p *MockPushRecord) Type() string {
	return p.appendType
}

func (p *MockPushRecord) ID() string {
	return p.id
}

func (p *MockPushRecord) IPAddr() string {
	return p.ipAddr
}

func (p *MockPushRecord) Text() string {
	return p.text
}

func (p *MockPushRecord) Time() time.Time {
	return p.time
}
