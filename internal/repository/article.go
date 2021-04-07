package repository

import (
	"context"

	"github.com/Ptt-official-app/go-bbs"
)

// PopularArticleRecord is an ArticleRecord which has boardID information.
type PopularArticleRecord interface {
	// Note: go-bbs has not implemented this yet
	// TODO: use bbs.PopularArticleRecord or something when it is ready
	bbs.ArticleRecord
	BoardID() string
}

func (repo *repository) GetPopularArticles(ctx context.Context) ([]PopularArticleRecord, error) {
	// Note: go-bbs has not implemented this yet
	// TODO: delegate to repo.db when it is ready
	return []PopularArticleRecord{}, nil
}

func (repo *repository) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (map[string]interface{}, error) {
	return nil, nil
}
