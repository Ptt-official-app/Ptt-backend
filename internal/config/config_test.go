package config

import (
	"bytes"
	"testing"
	"text/template"
	"time"
)

type MockUserRecord struct {
}

func (record *MockUserRecord) UserID() string {
	return "unknown"
}

func (record *MockUserRecord) Nickname() string {
	return "無名"
}

func (record *MockUserRecord) LastHost() string {
	return "172.17.0.3"
}

type MockArticleRecord struct {
}

func (article *MockArticleRecord) Title() string {
	return "[說明] this is a test article 這是一篇測試文章"
}

func (article *MockArticleRecord) Filename() string {
	return "M.1633153171.A.221"
}

func TestNewConfig(t *testing.T) {
	actual, err := NewConfig("testcase/01.toml", "")
	if err != nil {
		t.Errorf("NewConfig error excepted nil, got %v", err)
		return
	}

	expected := Config{
		BBSHome: "./home/bbs",
		// ListenPort            int16
		// AccessTokenPrivateKey string
		// AccessTokenPublicKey  string
		// AccessTokenExpiresAt  time.Duration
	}

	if actual.BBSHome != expected.BBSHome {
		t.Errorf("bbshome not match, expected: %v, got: %v", expected.BBSHome, actual.BBSHome)
		return
	}

	temp, err := template.New("Ptt-article-template").Parse(PttArticleTemplate)
	if err != nil {
		t.Fatal(err)
	}

	boardID := "test"
	content := "你好123 hi"
	currentTime := time.Now()
	buffer := bytes.NewBuffer(nil)
	err = temp.Execute(buffer, ArticleArguments{
		UserData:      &MockUserRecord{},
		Article:       &MockArticleRecord{},
		BoardID:       boardID,
		Content:       content,
		PostANSICDate: currentTime.Format(time.ANSIC),
	})

	if err != nil {
		t.Fatal(err)
	}
}

// the following mock function is not used... for now.
func (article *MockArticleRecord) Modified() time.Time {
	return time.Now()
}

func (article *MockArticleRecord) SetModified(newModified time.Time) {
}

func (article *MockArticleRecord) Recommend() int {
	return 0
}

func (article *MockArticleRecord) Date() string {
	return ""
}
func (article *MockArticleRecord) Money() int {
	return 0
}
func (article *MockArticleRecord) Owner() string {
	return "nothing"
}

func (record *MockUserRecord) HashedPassword() string {
	return ""
}
func (record *MockUserRecord) VerifyPassword(password string) error {
	return nil
}
func (record *MockUserRecord) RealName() string {
	return ""
}
func (record *MockUserRecord) NumLoginDays() int {
	return 0
}
func (record *MockUserRecord) NumPosts() int {
	return 0
}
func (record *MockUserRecord) Money() int {
	return 0
}
func (record *MockUserRecord) LastLogin() time.Time {
	return time.Now()
}
func (record *MockUserRecord) UserFlag() uint32 {
	return 0
}
