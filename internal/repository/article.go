package repository

import (
	"bytes"
	"container/heap"
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"text/template"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/go-bbs"
)

const (
	popularArticlesLimit    = 100
	popularArticlesCacheTTL = 5 * time.Minute
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

type popularArticlesCache struct {
	mu        sync.Mutex
	items     []PopularArticleRecord
	expiresAt time.Time
}

func clonePopularArticles(items []PopularArticleRecord) []PopularArticleRecord {
	result := make([]PopularArticleRecord, len(items))
	copy(result, items)
	return result
}

func (cache *popularArticlesCache) get(now time.Time, refresh func() ([]PopularArticleRecord, error)) ([]PopularArticleRecord, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.items != nil && now.Before(cache.expiresAt) {
		return clonePopularArticles(cache.items), nil
	}

	items, err := refresh()
	if err != nil {
		return nil, err
	}

	cache.items = clonePopularArticles(items)
	cache.expiresAt = now.Add(popularArticlesCacheTTL)
	return clonePopularArticles(cache.items), nil
}

func (cache *popularArticlesCache) invalidate() {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.items = nil
	cache.expiresAt = time.Time{}
}

func isMorePopular(a, b PopularArticleRecord) bool {
	if a.Recommend() != b.Recommend() {
		return a.Recommend() > b.Recommend()
	}
	if !a.Modified().Equal(b.Modified()) {
		return a.Modified().After(b.Modified())
	}
	if a.BoardID() != b.BoardID() {
		return a.BoardID() < b.BoardID()
	}
	return a.Filename() < b.Filename()
}

type popularArticleMinHeap []PopularArticleRecord

func (h popularArticleMinHeap) Len() int { return len(h) }
func (h popularArticleMinHeap) Less(i, j int) bool {
	// container/heap is a min-heap; keep the least popular candidate at root.
	return isMorePopular(h[j], h[i])
}
func (h popularArticleMinHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *popularArticleMinHeap) Push(value interface{}) {
	*h = append(*h, value.(PopularArticleRecord))
}
func (h *popularArticleMinHeap) Pop() interface{} {
	old := *h
	last := len(old) - 1
	item := old[last]
	old[last] = nil
	*h = old[:last]
	return item
}

type articleRecordsReader func(context.Context, string) ([]bbs.ArticleRecord, error)

func buildPopularArticles(ctx context.Context, boards []bbs.BoardRecord, readArticles articleRecordsReader) ([]PopularArticleRecord, error) {
	candidates := &popularArticleMinHeap{}
	heap.Init(candidates)

	for _, board := range boards {
		if !BoardIsPublic(board) {
			continue
		}

		articles, err := readArticles(ctx, board.BoardID())
		if err != nil {
			return nil, fmt.Errorf("read articles from board %s: %w", board.BoardID(), err)
		}

		for _, article := range articles {
			if article == nil {
				continue
			}
			candidate := &PopularArticle{
				filename:       article.Filename(),
				modified:       article.Modified(),
				recommendCount: article.Recommend(),
				owner:          article.Owner(),
				date:           article.Date(),
				title:          article.Title(),
				money:          article.Money(),
				boardID:        board.BoardID(),
			}

			if candidates.Len() < popularArticlesLimit {
				heap.Push(candidates, candidate)
				continue
			}
			if isMorePopular(candidate, (*candidates)[0]) {
				heap.Pop(candidates)
				heap.Push(candidates, candidate)
			}
		}
	}

	result := make([]PopularArticleRecord, candidates.Len())
	copy(result, *candidates)
	sort.Slice(result, func(i, j int) bool {
		return isMorePopular(result[i], result[j])
	})
	return result, nil
}

func (repo *repository) GetPopularArticles(ctx context.Context) ([]PopularArticleRecord, error) {
	return repo.popularArticles.get(time.Now(), func() ([]PopularArticleRecord, error) {
		return buildPopularArticles(ctx, repo.GetBoards(ctx), func(ctx context.Context, boardID string) ([]bbs.ArticleRecord, error) {
			return repo.GetBoardArticleRecords(ctx, boardID, 0, ^uint(0))
		})
	})
}

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

func (repo *repository) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (PushRecord, error) {
	// Append comment into board article file
	now := time.Now()
	result := "[1;31m " + appendType + " " + userID + " [m [33m: " + text + " [m " + now.Format("01/02 15:04") + "\n"
	err := repo.db.AppendBoardArticleFile(boardID, filename, bbs.Utf8ToBig5(result))
	if err != nil {
		return nil, err
	}
	repo.popularArticles.invalidate()

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
	repo.popularArticles.invalidate()

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
