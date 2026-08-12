package http

import "context"

func (usecase *MockUsecase) DeleteArticle(ctx context.Context, token, boardID, filename string) error {
	return nil
}
