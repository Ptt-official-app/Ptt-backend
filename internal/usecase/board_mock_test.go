package usecase

import (
	"context"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/go-bbs"
)

func (repo *MockRepository) GetBoards(ctx context.Context) []bbs.BoardRecord {
	return []bbs.BoardRecord{
		NewMockBoardRecord("SYSOP", "SYSOP", "嘰哩 ◎站長好!", false),
	}
}

func (repo *MockRepository) GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error) {
	return []byte{}, nil
}

func (repo *MockRepository) GetBoardArticleRecords(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error) {
	articleRecords := []*MockArticleRecord{
		{
			filename:       "filename1",
			modified:       time.Time{},
			recommendCount: 10,
			owner:          "SYSOP",
			date:           "",
			title:          "[討論] 偶爾要發個廢文",
			money:          0,
		},
		{
			filename:       "filename2",
			modified:       time.Time{},
			recommendCount: -20,
			owner:          "XDXD",
			date:           "",
			title:          "[問題] UNICODE",
			money:          0,
		},
	}
	result := make([]bbs.ArticleRecord, len(articleRecords))
	for i, v := range articleRecords {
		result[i] = v
	}
	return result, nil
}

func (repo *MockRepository) GetPopularArticles(ctx context.Context) ([]repository.PopularArticleRecord, error) {
	articleRecords := []*MockPopularArticle{
		{
			filename:       "",
			modified:       time.Time{},
			recommendCount: 10,
			owner:          "SYSOP",
			date:           "",
			title:          "Popular Article 1",
			money:          0,
			boardID:        "Gossiping",
		},
		{
			filename:       "",
			modified:       time.Time{},
			recommendCount: -20,
			owner:          "XDXD",
			date:           "",
			title:          "Popular Article 2",
			money:          0,
			boardID:        "Gossiping",
		},
		{
			filename:       "",
			modified:       time.Time{},
			recommendCount: 0,
			owner:          "TEST",
			date:           "",
			title:          "Popular Article 3",
			money:          0,
			boardID:        "Joke",
		},
	}
	result := make([]repository.PopularArticleRecord, len(articleRecords))
	for i, v := range articleRecords {
		result[i] = v
	}
	return result, nil
}

func (repo *MockRepository) CreateArticle(ctx context.Context, userID, boardID, title, article string) (bbs.ArticleRecord, error) {
	return  &MockArticleRecord{
		filename:       "filename1",
		modified:       time.Time{},
		recommendCount: 10,
		owner:          "SYSOP",
		date:           "",
		title:          "[討論] 偶爾要發個廢文",
		money:          0,
	}, nil
}

func (repo *MockRepository) GetRawArticle(boardID, filename string) (string, error) {
	return "test", nil
}

func (repo *MockRepository) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (repository.PushRecord, error) {
	return MockPushRecord{
		appendType: appendType,
		userID:     userID,
		text:       text,
		time:       time.Time{},
		ipAddr:     "127.0.0.1",
	}, nil
}

func (repo *MockRepository) ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (repository.ForwardArticleToBoardRecord, error) {
	return &MockForwardArticleToBoardRecord{
		filename:       "",
		modified:       time.Time{},
		recommendCount: 10,
		owner:          "rico",
		date:           "",
		title:          "[閒聊] 可不可以當 couch potato",
		money:          0,
		destBoardID:    "Soft_Job",
		ipAddr:         "1.1.1.1",
		forwardTime:    time.Time{},
		forwardTitle:   "轉 [閒聊] 可不可以當 couch potato",
	}, nil
}

func (repo *MockRepository) GetBoardTreasureRecords(ctx context.Context, boardID string, treasureIDs []string) ([]bbs.ArticleRecord, error) {
	return []bbs.ArticleRecord{}, nil
}

type MockBoardRecord struct {
	boardID          string
	title            string
	isClass          bool
	classID          string
	postLimitPosts   uint8
	postLimitLogins  uint8
	postLimitBadPost uint8
}

func NewMockBoardRecord(classID, boardID, title string, isClass bool) *MockBoardRecord {
	return &MockBoardRecord{boardID: boardID, title: title, isClass: isClass, classID: classID}
}

func (b *MockBoardRecord) BoardID() string            { return b.boardID }
func (b *MockBoardRecord) Title() string              { return b.title }
func (b *MockBoardRecord) IsClass() bool              { return b.isClass }
func (b *MockBoardRecord) ClassID() string            { return b.classID }
func (b *MockBoardRecord) BM() []string               { return make([]string, 0) }
func (b *MockBoardRecord) GetPostLimitPosts() uint8   { return b.postLimitPosts }
func (b *MockBoardRecord) GetPostLimitLogins() uint8  { return b.postLimitLogins }
func (b *MockBoardRecord) GetPostLimitBadPost() uint8 { return b.postLimitBadPost }

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

type MockPopularArticle struct {
	filename       string
	modified       time.Time
	recommendCount int
	owner          string
	date           string
	title          string
	money          int
	boardID        string
}

func (m *MockPopularArticle) Filename() string               { return m.filename }
func (m *MockPopularArticle) Modified() time.Time            { return m.modified }
func (m *MockPopularArticle) SetModified(newValue time.Time) { m.modified = newValue }
func (m *MockPopularArticle) Recommend() int                 { return m.recommendCount }
func (m *MockPopularArticle) Date() string                   { return m.date }
func (m *MockPopularArticle) Title() string                  { return m.title }
func (m *MockPopularArticle) Money() int                     { return m.money }
func (m *MockPopularArticle) Owner() string                  { return m.owner }
func (m *MockPopularArticle) BoardID() string                { return m.boardID }

func (repo *MockRepository) GetUserArticles(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error) {
	articleRecords := []*MockArticleRecord{
		{
			filename:       "",
			modified:       time.Time{},
			recommendCount: 10,
			owner:          "user",
			date:           "",
			title:          "[討論] 薪水太少",
			money:          0,
		},
		{
			filename:       "",
			modified:       time.Time{},
			recommendCount: -20,
			owner:          "9487",
			date:           "",
			title:          "[問題] 我不會寫程式",
			money:          0,
		},
	}
	result := make([]bbs.ArticleRecord, len(articleRecords))
	for i, v := range articleRecords {
		result[i] = v
	}
	return result, nil
}

type MockPushRecord struct {
	appendType string
	userID     string
	ipAddr     string
	text       string
	time       time.Time
}

func (m MockPushRecord) Type() string    { return m.appendType }
func (m MockPushRecord) ID() string      { return m.userID }
func (m MockPushRecord) IPAddr() string  { return m.ipAddr }
func (m MockPushRecord) Text() string    { return m.text }
func (m MockPushRecord) Time() time.Time { return m.time }

type MockForwardArticleToBoardRecord struct {
	filename       string
	modified       time.Time
	recommendCount int
	owner          string
	date           string
	title          string
	money          int
	destBoardID    string
	ipAddr         string
	forwardTime    time.Time
	forwardTitle   string
}

func (m *MockForwardArticleToBoardRecord) Filename() string               { return m.filename }
func (m *MockForwardArticleToBoardRecord) Modified() time.Time            { return m.modified }
func (m *MockForwardArticleToBoardRecord) SetModified(newValue time.Time) { m.modified = newValue }
func (m *MockForwardArticleToBoardRecord) Recommend() int                 { return m.recommendCount }
func (m *MockForwardArticleToBoardRecord) Date() string                   { return m.date }
func (m *MockForwardArticleToBoardRecord) Title() string                  { return m.title }
func (m *MockForwardArticleToBoardRecord) Money() int                     { return m.money }
func (m *MockForwardArticleToBoardRecord) Owner() string                  { return m.owner }
func (m *MockForwardArticleToBoardRecord) DestBoardID() string            { return m.destBoardID }
func (m *MockForwardArticleToBoardRecord) IPAddr() string                 { return m.ipAddr }
func (m *MockForwardArticleToBoardRecord) ForwardTime() time.Time         { return m.forwardTime }
func (m *MockForwardArticleToBoardRecord) ForwardTitle() string           { return m.forwardTitle }
