package usecase

import (
	"context"
	"testing"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
)

func TestGetPopularArticles(t *testing.T) {
	resp := &MockRepository{}

	usecase := NewUsecase(&config.Config{}, resp)
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
