package usecase

import (
	"context"
	"testing"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
)

func TestGetPopularArticles(t *testing.T) {
	repo := &MockRepository{}

	usecase := NewUsecase(&config.Config{}, repo, &logging.DummyLogger{})
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

func TestForwardArticleToEmail(t *testing.T) {
	repo := &MockRepository{}

	userID := "mockUserID"
	boardID := "board1"
	filename := "filename1"
	email := "test@gmail.com"
	mail := &MockMail{}

	usecase := NewUsecase(&config.Config{}, repo, &logging.DummyLogger{})
	_ = usecase.UpdateMail(mail)
	err := usecase.ForwardArticleToEmail(context.TODO(), userID, boardID, filename, email)

	if err != nil {
		t.Errorf("ForwardArticleToEmail failed %v", err)
	}

	if mail.data["email"] != email {
		t.Errorf("Send Email with incorrect email, want %s, get %s\n", email, mail.data["email"])
	}

	if mail.data["title"] != "[討論] 偶爾要發個廢文" {
		t.Errorf("Send Email with incorrect title, want %s, get %s\n", "[討論] 偶爾要發個廢文", mail.data["title"])
	}

	if mail.data["userID"] != userID {
		t.Errorf("Send Email with incorrect userID, want %s, get %s\n", userID, mail.data["userID"])
	}
}

type MockMail struct {
	data map[string]interface{}
}

func (mail *MockMail) Send(email, title, userID string, body []byte) error {
	mail.data = map[string]interface{}{
		"email":  email,
		"title":  title,
		"userID": userID,
		"body":   body,
	}
	return nil
}
