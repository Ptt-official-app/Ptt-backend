package main

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/proto/api"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
	"github.com/Ptt-official-app/go-bbs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type searchBoard struct {
	bbs.UnimplementedBoardRecord
	id string
}

func (b searchBoard) BoardID() string { return b.id }

type searchArticle struct {
	bbs.UnimplementedArticleRecord
	filename, date, owner, title string
	recommend                    int
	modified                     time.Time
}

func (a searchArticle) Filename() string    { return a.filename }
func (a searchArticle) Date() string        { return a.date }
func (a searchArticle) Owner() string       { return a.owner }
func (a searchArticle) Title() string       { return a.title }
func (a searchArticle) Recommend() int      { return a.recommend }
func (a searchArticle) Modified() time.Time { return a.modified }

type searchUsecase struct {
	usecase.Usecase
	boards   []bbs.BoardRecord
	articles []bbs.ArticleRecord
}

func (u *searchUsecase) GetBoards(context.Context, string) []bbs.BoardRecord { return u.boards }

func (u *searchUsecase) GetBoardArticles(_ context.Context, _ string, _ uint, _ uint, cond *usecase.ArticleSearchCond) []bbs.ArticleRecord {
	if cond == nil {
		return u.articles
	}
	result := make([]bbs.ArticleRecord, 0, len(u.articles))
	for _, article := range u.articles {
		if cond.Title != "" && !strings.Contains(strings.ToLower(article.Title()), strings.ToLower(cond.Title)) {
			continue
		}
		if cond.ExactTitle != "" && article.Title() != cond.ExactTitle {
			continue
		}
		if cond.Author != "" && !strings.Contains(strings.ToLower(article.Owner()), strings.ToLower(cond.Author)) {
			continue
		}
		if cond.RecommendCountGreaterEqualIsSet && article.Recommend() < cond.RecommendCountGreaterEqual {
			continue
		}
		if cond.RecommendCountLessEqualIsSet && article.Recommend() > cond.RecommendCountLessEqual {
			continue
		}
		result = append(result, article)
	}
	return result
}

func newSearchServer() *server {
	return &server{usecase: &searchUsecase{
		boards: []bbs.BoardRecord{&searchBoard{id: "Test"}},
		articles: []bbs.ArticleRecord{
			&searchArticle{filename: "1", title: "Hello world", owner: "Alice", recommend: 10},
			&searchArticle{filename: "2", title: "Hello Go", owner: "Bob", recommend: -3},
			&searchArticle{filename: "3", title: "Other", owner: "Alice", recommend: 1},
		},
	}}
}

func searchRequest(filters ...*api.SearchFilter) *api.SearchRequest {
	return &api.SearchRequest{
		Ref:    &api.BoardRef{Ref: &api.BoardRef_Name{Name: "test"}},
		Filter: filters,
		Length: 20,
	}
}

func TestSearchFilters(t *testing.T) {
	cases := []struct {
		name   string
		filter *api.SearchFilter
		want   []string
	}{
		{"title", &api.SearchFilter{Type: api.SearchFilter_TYPE_TITLE, StringData: "hello"}, []string{"1", "2"}},
		{"exact title", &api.SearchFilter{Type: api.SearchFilter_TYPE_EXACT_TITLE, StringData: "Hello Go"}, []string{"2"}},
		{"author", &api.SearchFilter{Type: api.SearchFilter_TYPE_AUTHOR, StringData: "ali"}, []string{"1", "3"}},
		{"recommend >=", &api.SearchFilter{Type: api.SearchFilter_TYPE_RECOMMEND, NumberData: 5}, []string{"1"}},
		{"recommend <=", &api.SearchFilter{Type: api.SearchFilter_TYPE_RECOMMEND, NumberData: -2}, []string{"2"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reply, err := newSearchServer().Search(context.Background(), searchRequest(tc.filter))
			if err != nil {
				t.Fatal(err)
			}
			if got := postFilenames(reply.Posts); !equalStrings(got, tc.want) {
				t.Fatalf("filenames = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSearchEmptyAndPagination(t *testing.T) {
	s := newSearchServer()
	reply, err := s.Search(context.Background(), searchRequest(&api.SearchFilter{
		Type: api.SearchFilter_TYPE_TITLE, StringData: "missing",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if len(reply.Posts) != 0 || reply.TotalPosts != 0 {
		t.Fatalf("empty search = posts %d total %d", len(reply.Posts), reply.TotalPosts)
	}

	req := searchRequest()
	req.Offset = -2
	req.Length = 2
	reply, err = s.Search(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if got := postFilenames(reply.Posts); !equalStrings(got, []string{"2", "3"}) || reply.TotalPosts != 3 {
		t.Fatalf("negative page = %v total %d", postFilenames(reply.Posts), reply.TotalPosts)
	}

	req.Offset = 100
	reply, err = s.Search(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if len(reply.Posts) != 0 || reply.TotalPosts != 3 {
		t.Fatalf("offset past end = posts %d total %d", len(reply.Posts), reply.TotalPosts)
	}
}

func TestSearchRejectsInvalidRefsAndUnsupportedFilters(t *testing.T) {
	cases := []struct {
		name string
		req  *api.SearchRequest
		code codes.Code
	}{
		{"missing ref", &api.SearchRequest{Length: 1}, codes.InvalidArgument},
		{"bad board", &api.SearchRequest{Ref: &api.BoardRef{Ref: &api.BoardRef_Name{Name: "missing"}}, Length: 1}, codes.InvalidArgument},
		{"bad bid", &api.SearchRequest{Ref: &api.BoardRef{Ref: &api.BoardRef_Bid{Bid: 2}}, Length: 1}, codes.InvalidArgument},
		{"negative length", &api.SearchRequest{Ref: &api.BoardRef{Ref: &api.BoardRef_Name{Name: "test"}}, Length: -1}, codes.InvalidArgument},
		{"money unsupported", searchRequest(&api.SearchFilter{Type: api.SearchFilter_TYPE_MONEY}), codes.Unimplemented},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newSearchServer().Search(context.Background(), tc.req)
			if status.Code(err) != tc.code {
				t.Fatalf("status = %v, want %v; error = %v", status.Code(err), tc.code, err)
			}
		})
	}
}

func postFilenames(posts []*api.Post) []string {
	result := make([]string, len(posts))
	for i, post := range posts {
		result[i] = post.Filename
	}
	return result
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
