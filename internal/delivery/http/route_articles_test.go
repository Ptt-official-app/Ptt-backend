package http

import (
	"context"

	"github.com/PichuChen/go-bbs"
)

func (usecase *MockUsecase) GetPopularArticles(ctx context.Context) ([]bbs.ArticleRecord, error) {
	return []bbs.ArticleRecord{}, nil
}
