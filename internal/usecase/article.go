package usecase

import (
	"context"
	"fmt"

	"github.com/PichuChen/go-bbs"
)

func (usecase *usecase) GetPopularArticles(ctx context.Context) ([]bbs.ArticleRecord, error) {
	// articles should be sorted
	articles, err := usecase.repo.GetPopularArticles(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetPopularArticles error: %w", err)
	}
	return articles, nil
}
