package http

import (
	"context"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
	"github.com/Ptt-official-app/go-bbs"
)

// GetBoardByID returns the mock board record corresponding to boardID
func (usecase *MockUsecase) GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error) {
	boardRecord := NewMockBoardRecord("SYSOP", boardID, "嘰哩 ◎站長好!", false)
	return boardRecord, nil
}

// GetBoards returns the mock board records available for specific users identified by userID
func (usecase *MockUsecase) GetBoards(ctx context.Context, userID string) []bbs.BoardRecord {
	result := make([]bbs.BoardRecord, 0)
	result = append(result, NewMockBoardRecord("junk", "TEST", "發電 ◎雜七雜八的垃圾", false))
	return result
}

// GetPopularBoards returns the mock popular boards
func (usecase *MockUsecase) GetPopularBoards(ctx context.Context) ([]bbs.BoardRecord, error) {
	result := make([]bbs.BoardRecord, 0)
	result = append(result, NewMockBoardRecord("SYSOP", "", "嘰哩 ◎站長好!", true))
	result = append(result, NewMockBoardRecord("junk", "TEST", "發電 ◎雜七雜八的垃圾", false))
	return result, nil
}

// GetBoardArticles returns the mock board articles.
func (usecase *MockUsecase) GetBoardArticles(ctx context.Context, boardID string, cond *usecase.ArticleSearchCond) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"filename":        "test-001",
			"modified_time":   "2009-01-01T12:59:59Z",
			"recommend_count": 9,
			"post_date":       "2009-01-01",
			"title":           "post for testing",
			"money":           "10",
			"owner":           "tester",
			"url":             "http://test/test-001.html",
		},
	}
}

// TODO: Implement GetBoardArticle
func (usecase *MockUsecase) GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error) {
	panic("Not implemented")
}

// GetBoardTreasures returns mock treasures for specific board identified by boardID and treasuresID
func (usecase *MockUsecase) GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"filename":  "testBoardTreasures",
			"post_date": "2020-03-12",
			"title":     "testing",
			"owner":     "ptt",
			"url":       "https://google.com",
		},
	}
}

// GetBoardPostsLimitation returns mock board post limitation with posts limit, logins limit, and bad post limit are all equal to 0
func (usecase *MockUsecase) GetBoardPostsLimitation(ctx context.Context, boardID string) (*usecase.BoardPostLimitation, error) {
	return NewMockBoardPostLimitation(0, 0, 0), nil
}

func NewMockBoardPostLimitation(postsLimit uint8, loginsLimit uint8, badPostLimit uint8) *usecase.BoardPostLimitation {
	return &usecase.BoardPostLimitation{
		PostsLimit:   postsLimit,
		LoginsLimit:  loginsLimit,
		BadPostLimit: badPostLimit,
	}
}

type MockBoardRecord struct {
	boardID string
	title   string
	isClass bool
	classID string
}

func NewMockBoardRecord(classID, boardID, title string, isClass bool) *MockBoardRecord {
	return &MockBoardRecord{boardID: boardID, title: title, isClass: isClass, classID: classID}
}

func (b *MockBoardRecord) BoardID() string          { return b.boardID }
func (b *MockBoardRecord) Title() string            { return b.title }
func (b *MockBoardRecord) IsClass() bool            { return b.isClass }
func (b *MockBoardRecord) ClassID() string          { return b.classID }
func (b *MockBoardRecord) BM() []string             { return make([]string, 0) }
func (b *MockBoardRecord) IsNoCount() bool          { return false }
func (b *MockBoardRecord) IsGroupBoard() bool       { return false }
func (b *MockBoardRecord) IsHide() bool             { return false }
func (b *MockBoardRecord) IsPostMask() bool         { return false }
func (b *MockBoardRecord) IsAnonymous() bool        { return false }
func (b *MockBoardRecord) IsDefaultAnonymous() bool { return false }
func (b *MockBoardRecord) IsNoCredit() bool         { return false }
func (b *MockBoardRecord) IsVoteBoard() bool        { return false }
func (b *MockBoardRecord) IsWarnEL() bool           { return false }
func (b *MockBoardRecord) IsTop() bool              { return false }
func (b *MockBoardRecord) IsNoRecommend() bool      { return false }
func (b *MockBoardRecord) IsAngelAnonymous() bool   { return false }
func (b *MockBoardRecord) IsBMCount() bool          { return false }
func (b *MockBoardRecord) IsIsSymbolic() bool       { return false }
func (b *MockBoardRecord) IsNoBoo() bool            { return false }
func (b *MockBoardRecord) IsRestrictedPost() bool   { return false }
func (b *MockBoardRecord) IsGuestPost() bool        { return false }
func (b *MockBoardRecord) IsCooldown() bool         { return false }
func (b *MockBoardRecord) IsCPLog() bool            { return false }
func (b *MockBoardRecord) IsNoFastRecommend() bool  { return false }
func (b *MockBoardRecord) IsIPLogRecommend() bool   { return false }
func (b *MockBoardRecord) IsOver18() bool           { return false }
func (b *MockBoardRecord) IsNoReply() bool          { return false }
func (b *MockBoardRecord) IsAlignedComment() bool   { return false }
func (b *MockBoardRecord) IsNoSelfDeletePost() bool { return false }
func (b *MockBoardRecord) IsBMMaskContent() bool    { return false }
