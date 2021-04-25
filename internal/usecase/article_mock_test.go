package usecase

import (
	"context"
	"time"
)

func (repo *MockRepository) ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"user_id":      "rico",
		"article_id":   "12345",
		"forward_time": time.Time{},
		"title":        "[閒聊] 可不可以當 couch potato",
		"dest_board":   "Gossiping",
	}
	return result, nil
}

func (repo *MockRepository) ForwardArticleToEmail(ctx context.Context, userID, boardID, filename, email string) error {
	return nil
}
