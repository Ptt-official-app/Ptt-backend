package usecase

import (
	"context"
	"testing"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/go-bbs"
)

func TestSearchArticles(t *testing.T) {
	repository := &MockRepository{}
	articleRecords, _ := repository.GetBoardArticleRecords(context.Background(), "")

	type TestCase struct {
		input         *ArticleSearchCond
		expectedItems []bbs.ArticleRecord
	}

	mockArticle1 := MockArticleRecord{
		recommendCount: 10,
		owner:          "SYSOP",
		title:          "[討論] 偶爾要發個廢文",
	}

	mockArticle2 := MockArticleRecord{
		recommendCount: -20,
		owner:          "XDXD",
		title:          "[問題] UNICODE",
	}

	cases := []TestCase{
		{
			input: &ArticleSearchCond{
				Title:                           "廢文",
				Author:                          "s",
				RecommendCountGreaterEqual:      0,
				RecommendCountLessEqual:         50,
				RecommendCountGreaterEqualIsSet: true,
				RecommendCountLessEqualIsSet:    true,
			},
			expectedItems: []bbs.ArticleRecord{&mockArticle1},
		},
		{
			input: &ArticleSearchCond{
				Title:                        "",
				Author:                       "X",
				RecommendCountLessEqual:      10,
				RecommendCountLessEqualIsSet: true,
			},
			expectedItems: []bbs.ArticleRecord{&mockArticle2},
		},
		{
			input: &ArticleSearchCond{
				Title: "Unicode",
			},
			expectedItems: []bbs.ArticleRecord{&mockArticle2},
		},
		{
			input: &ArticleSearchCond{
				Author: "sysop",
			},
			expectedItems: []bbs.ArticleRecord{&mockArticle1},
		},
	}

	for index, c := range cases {
		cond := c.input
		expectedItems := c.expectedItems

		actualItems := searchArticles(articleRecords, cond)
		for i, v := range actualItems {
			if v.Title() != expectedItems[i].Title() ||
				v.Owner() != expectedItems[i].Owner() ||
				v.Recommend() != expectedItems[i].Recommend() {
				t.Errorf("article not match on index %d, expected: %v, got: %v", index, expectedItems[i], v)
			}
		}
	}
}

func TestGetBoardPostsLimitation(t *testing.T) {
	resp := &MockRepository{}

	usecase := NewUsecase(&config.Config{}, resp, logging.DefaultDummyLogger)

	limitation, err := usecase.GetBoardPostsLimitation(context.TODO(), "board-id")
	if err != nil {
		t.Errorf("getBoardPostsLimitation with board-id excepted not nil error, got nil")
		return
	}

	if limitation.PostsLimit != 0 {
		t.Errorf("limitation.PostsLimit is excepted 0, got %d", limitation.PostsLimit)
		return
	}

	if limitation.LoginsLimit != 0 {
		t.Errorf("limitation.LoginsLimit is excepted 0, got %d", limitation.LoginsLimit)
		return
	}

	if limitation.BadPostLimit != 0 {
		t.Errorf("limitation.BadPostLimit is excepted 0, got %d", limitation.BadPostLimit)
		return
	}
}

func TestGetPopularBoards(t *testing.T) {
	resp := &MockRepository{}
	usecase := NewUsecase(&config.Config{}, resp, logging.DefaultDummyLogger)

	filtedBoards, err := usecase.GetPopularBoards(context.TODO())
	if err != nil {
		t.Errorf("GetPopularBoards error : %v", err)
	}

	// TODO:The test here must be changed after the implementation is completed
	if len(filtedBoards) != 0 {
		t.Errorf("GetPopularBoards shoule have no data")
	}
}
