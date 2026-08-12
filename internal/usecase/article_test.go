package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/go-bbs"
)

func TestGetPopularArticles(t *testing.T) {
	repo := &MockRepository{}

	usecase := NewUsecase(&config.Config{}, repo)
	articles, err := usecase.GetPopularArticles(context.TODO())
	if err != nil {
		t.Errorf("GetPopularArticles expected err == nil, got %v", err)
	}

	if len(articles) != 3 {
		t.Errorf("GetPopularArticles should return 3 articles, got %v", len(articles))
	}

	expectedFirstArticleTitle := "Popular Article 1"
	if articles[0].Title() != expectedFirstArticleTitle {
		t.Errorf("GetPopularArticles should return first article with title %s, got %s", expectedFirstArticleTitle, articles[0].Title())
	}
}

type usefulnessTestRepository struct {
	*MockRepository
	article    []byte
	owner      string
	filename   string
	appendType []string
}

func newUsefulnessTestRepository(owner, filename, article string) *usefulnessTestRepository {
	return &usefulnessTestRepository{
		MockRepository: &MockRepository{},
		article:        bbs.Utf8ToBig5(article),
		owner:          owner,
		filename:       filename,
	}
}

func (repo *usefulnessTestRepository) GetBoardArticleRecords(context.Context, string, uint, uint) ([]bbs.ArticleRecord, error) {
	return []bbs.ArticleRecord{
		&MockArticleRecord{
			filename: repo.filename,
			owner:    repo.owner,
		},
	}, nil
}

func (repo *usefulnessTestRepository) GetBoardArticle(context.Context, string, string) ([]byte, error) {
	result := make([]byte, len(repo.article))
	copy(result, repo.article)
	return result, nil
}

func (repo *usefulnessTestRepository) AppendComment(_ context.Context, userID, _, _ string, appendType, text string) (repository.PushRecord, error) {
	line := "[1;31m " + appendType + " " + userID + " [m [33m: " + text + " [m 01/02 03:04\n"
	repo.article = append(repo.article, bbs.Utf8ToBig5(line)...)
	repo.appendType = append(repo.appendType, appendType)
	return MockPushRecord{
		appendType: appendType,
		userID:     userID,
		text:       line,
		time:       time.Unix(0, 0),
	}, nil
}

func TestUpdateUsefulnessStateTransitions(t *testing.T) {
	const (
		owner    = "user01"
		userID   = "user02"
		boardID  = "board1"
		filename = "filename1"
	)

	repo := newUsefulnessTestRepository(owner, filename, "文章本文\n")
	uc := NewUsecase(&config.Config{}, repo)

	assertScore := func(want int) {
		t.Helper()
		got := userUsefulnessScore(bbs.Big5ToUtf8(repo.article), userID)
		if got != want {
			t.Fatalf("usefulness score = %d, want %d", got, want)
		}
	}

	if _, err := uc.UpdateUsefulness(context.Background(), userID, boardID, filename, "↑"); err != nil {
		t.Fatalf("first up vote failed: %v", err)
	}
	assertScore(1)

	if _, err := uc.UpdateUsefulness(context.Background(), userID, boardID, filename, "↑"); err == nil {
		t.Fatal("second up vote should be rejected at +1")
	}
	assertScore(1)

	if _, err := uc.UpdateUsefulness(context.Background(), userID, boardID, filename, "↓"); err != nil {
		t.Fatalf("down vote from +1 to 0 failed: %v", err)
	}
	assertScore(0)

	if _, err := uc.UpdateUsefulness(context.Background(), userID, boardID, filename, "↓"); err != nil {
		t.Fatalf("down vote from 0 to -1 failed: %v", err)
	}
	assertScore(-1)

	if _, err := uc.UpdateUsefulness(context.Background(), userID, boardID, filename, "↓"); err == nil {
		t.Fatal("third down vote should be rejected at -1")
	}
	assertScore(-1)

	wantAppended := []string{"↑", "↓", "↓"}
	if len(repo.appendType) != len(wantAppended) {
		t.Fatalf("appended %v, want %v", repo.appendType, wantAppended)
	}
	for i := range wantAppended {
		if repo.appendType[i] != wantAppended[i] {
			t.Fatalf("appendType[%d] = %q, want %q", i, repo.appendType[i], wantAppended[i])
		}
	}
}

func TestUpdateUsefulnessRejectsOwnArticle(t *testing.T) {
	repo := newUsefulnessTestRepository("User01", "filename1", "文章本文\n")
	uc := NewUsecase(&config.Config{}, repo)

	if _, err := uc.UpdateUsefulness(context.Background(), "user01", "board1", "filename1", "↑"); err == nil {
		t.Fatal("article owner should not be able to rate own article")
	}
	if len(repo.appendType) != 0 {
		t.Fatalf("own-article vote wrote comments: %v", repo.appendType)
	}
}

func TestUpdateUsefulnessRejectsUnknownArticleAndType(t *testing.T) {
	repo := newUsefulnessTestRepository("user01", "filename1", "文章本文\n")
	uc := NewUsecase(&config.Config{}, repo)

	if _, err := uc.UpdateUsefulness(context.Background(), "user02", "board1", "missing", "↑"); err == nil {
		t.Fatal("unknown article should be rejected")
	}
	if _, err := uc.UpdateUsefulness(context.Background(), "user02", "board1", "filename1", "推"); err == nil {
		t.Fatal("non-arrow usefulness type should be rejected")
	}
}

func TestUserUsefulnessScoreParsesArrowComments(t *testing.T) {
	article := stringsJoinLines(
		"一般本文提到 ↑ user02 但不是推文格式",
		"[1;31m ↑ user02 [m [33m: [m 01/02 03:04",
		"\x1b[1;31m↓ USER02 \x1b[m: \x1b[33m \x1b[m 01/02 03:05",
		"↑ otherUser :",
		"↓ user02 :",
	)

	if got := userUsefulnessScore(article, "user02"); got != -1 {
		t.Fatalf("userUsefulnessScore() = %d, want -1", got)
	}
}

func stringsJoinLines(lines ...string) string {
	result := ""
	for _, line := range lines {
		result += line + "\n"
	}
	return result
}

func TestUpdateUsefulness(t *testing.T) {
	repo := &MockRepository{}
	userID := "mockUserID"
	boardID := "board1"
	filename := "filename1"
	appendType := "↑"

	usecase := NewUsecase(&config.Config{}, repo)

	record, err := usecase.UpdateUsefulness(context.TODO(), userID, boardID, filename, appendType)

	if err != nil {
		t.Errorf("UpdateUsefulness expected err == nil, got %v", err)
	}

	if record.Type() != appendType {
		t.Errorf("Push record with incorrect appendType, want %s, get %s", appendType, record.Type())
	}

	if record.ID() != userID {
		t.Errorf("Push record with incorrect userID, want %s, get %s", userID, record.ID())
	}
}

func TestForwardArticleToEmail(t *testing.T) {
	repo := &MockRepository{}

	userID := "mockUserID"
	boardID := "board1"
	filename := "filename1"
	email := "test@gmail.com"
	mail := &MockMail{}

	usecase := NewUsecase(&config.Config{}, repo)
	_ = usecase.UpdateMail(mail)
	err := usecase.ForwardArticleToEmail(context.TODO(), userID, boardID, filename, email)
	if err != nil {
		t.Errorf("ForwardArticleToEmail failed %v", err)
	}

	if mail.data["email"] != email {
		t.Errorf("Send Email with incorrect email, want %s, get %s\n", email, mail.data["email"])
	}

	if mail.data["title"] != "[討論] 偶爾要發個廢文" {
		t.Errorf("Send Email with incorrect title, want %s, get %s\n", "[討論] 偶爾要發個廢文", mail.data["title"])
	}

	if mail.data["userID"] != userID {
		t.Errorf("Send Email with incorrect userID, want %s, get %s\n", userID, mail.data["userID"])
	}
}

type MockMail struct {
	data map[string]interface{}
}

func (mail *MockMail) Send(email, title, userID string, body []byte) error {
	mail.data = map[string]interface{}{
		"email":  email,
		"title":  title,
		"userID": userID,
		"body":   body,
	}
	return nil
}
