package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-bbs"
	"github.com/Ptt-official-app/go-bbs/pttbbs"
)

func TestBuildPopularArticlesFiltersRanksAndLimits(t *testing.T) {
	publicBoard := &pttbbs.BoardHeader{BrdName: "Public"}
	restrictedBoard := &pttbbs.BoardHeader{BrdName: "Restricted", Level: pttbbs.PermBM}

	publicArticles := make([]bbs.ArticleRecord, 0, 105)
	for i := 0; i < 105; i++ {
		publicArticles = append(publicArticles, &PopularArticle{
			filename:       fmt.Sprintf("M.%03d.A.000", i),
			modified:       time.Unix(int64(i), 0),
			recommendCount: i,
			owner:          "user",
			title:          fmt.Sprintf("article %d", i),
		})
	}

	articlesByBoard := map[string][]bbs.ArticleRecord{
		"Public": publicArticles,
		"Restricted": {
			&PopularArticle{
				filename:       "M.999.A.000",
				modified:       time.Unix(999, 0),
				recommendCount: 999,
				owner:          "sysop",
				title:          "must not leak",
			},
		},
	}
	readBoards := make([]string, 0)

	result, err := buildPopularArticles(
		context.Background(),
		[]bbs.BoardRecord{publicBoard, restrictedBoard},
		func(_ context.Context, boardID string) ([]bbs.ArticleRecord, error) {
			readBoards = append(readBoards, boardID)
			return articlesByBoard[boardID], nil
		},
	)
	if err != nil {
		t.Fatalf("buildPopularArticles() error = %v", err)
	}

	if len(readBoards) != 1 || readBoards[0] != "Public" {
		t.Fatalf("read boards = %v, want only Public", readBoards)
	}
	if len(result) != popularArticlesLimit {
		t.Fatalf("len(result) = %d, want %d", len(result), popularArticlesLimit)
	}
	if result[0].Recommend() != 104 {
		t.Fatalf("highest recommend = %d, want 104", result[0].Recommend())
	}
	if result[len(result)-1].Recommend() != 5 {
		t.Fatalf("lowest retained recommend = %d, want 5", result[len(result)-1].Recommend())
	}
	for i := 1; i < len(result); i++ {
		if result[i-1].Recommend() < result[i].Recommend() {
			t.Fatalf("result is not sorted at %d: %d < %d", i, result[i-1].Recommend(), result[i].Recommend())
		}
		if result[i].BoardID() != "Public" {
			t.Fatalf("restricted board leaked into result: %s", result[i].BoardID())
		}
	}
}

func TestBuildPopularArticlesUsesDeterministicTieBreak(t *testing.T) {
	board := &pttbbs.BoardHeader{BrdName: "Public"}
	modified := time.Unix(100, 0)
	articles := []bbs.ArticleRecord{
		&PopularArticle{filename: "M.2.A.000", modified: modified, recommendCount: 10},
		&PopularArticle{filename: "M.1.A.000", modified: modified, recommendCount: 10},
		&PopularArticle{filename: "M.3.A.000", modified: modified.Add(time.Second), recommendCount: 10},
	}

	result, err := buildPopularArticles(
		context.Background(),
		[]bbs.BoardRecord{board},
		func(_ context.Context, _ string) ([]bbs.ArticleRecord, error) {
			return articles, nil
		},
	)
	if err != nil {
		t.Fatalf("buildPopularArticles() error = %v", err)
	}

	want := []string{"M.3.A.000", "M.1.A.000", "M.2.A.000"}
	for i, filename := range want {
		if result[i].Filename() != filename {
			t.Fatalf("result[%d].Filename() = %s, want %s", i, result[i].Filename(), filename)
		}
	}
}

func TestPopularArticlesCache(t *testing.T) {
	cache := &popularArticlesCache{}
	now := time.Unix(1000, 0)
	refreshCalls := 0
	refresh := func() ([]PopularArticleRecord, error) {
		refreshCalls++
		return []PopularArticleRecord{
			&PopularArticle{filename: fmt.Sprintf("M.%d.A.000", refreshCalls)},
		}, nil
	}

	first, err := cache.get(now, refresh)
	if err != nil {
		t.Fatalf("first cache get error = %v", err)
	}
	first[0] = nil

	second, err := cache.get(now.Add(popularArticlesCacheTTL-time.Second), refresh)
	if err != nil {
		t.Fatalf("second cache get error = %v", err)
	}
	if refreshCalls != 1 {
		t.Fatalf("refresh calls = %d before expiry, want 1", refreshCalls)
	}
	if second[0] == nil || second[0].Filename() != "M.1.A.000" {
		t.Fatalf("cached result was mutated through returned slice")
	}

	_, err = cache.get(now.Add(popularArticlesCacheTTL), refresh)
	if err != nil {
		t.Fatalf("expired cache get error = %v", err)
	}
	if refreshCalls != 2 {
		t.Fatalf("refresh calls = %d after expiry, want 2", refreshCalls)
	}

	cache.invalidate()
	_, err = cache.get(now.Add(popularArticlesCacheTTL+time.Second), refresh)
	if err != nil {
		t.Fatalf("invalidated cache get error = %v", err)
	}
	if refreshCalls != 3 {
		t.Fatalf("refresh calls = %d after invalidate, want 3", refreshCalls)
	}
}
