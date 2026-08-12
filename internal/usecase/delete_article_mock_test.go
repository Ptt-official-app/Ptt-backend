package usecase

import "context"

func (repo *MockRepository) DeleteBoardArticle(ctx context.Context, boardID, filename, deletedBy string) error {
	return nil
}
