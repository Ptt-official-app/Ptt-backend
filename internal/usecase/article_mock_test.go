package usecase

import (
	"context"
	"time"

	"github.com/PichuChen/go-bbs"
)

func (repo *MockRepository) GetPopularArticles(ctx context.Context) ([]bbs.ArticleRecord, error) {
	ret := []bbs.ArticleRecord{
		&MockArticle{"Popular Article 1"},
		&MockArticle{"Popular Article 2"},
		&MockArticle{"Popular Article 3"},
	}
	return ret, nil
}

type MockArticle struct {
	title string
}

func (a *MockArticle) Filename() string    { return "" }
func (a *MockArticle) Modified() time.Time { return time.Unix(0, 0) }
func (a *MockArticle) Recommend() int      { return 0 }
func (a *MockArticle) Date() string        { return "" }
func (a *MockArticle) Title() string       { return a.title }
func (a *MockArticle) Money() int          { return 0 }
func (a *MockArticle) Owner() string       { return "" }
