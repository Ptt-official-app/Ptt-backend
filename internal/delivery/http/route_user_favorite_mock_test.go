package http

import "context"

func (usecase *MockUsecase) AddUserFavorite(ctx context.Context, token, userID, favoriteType, boardID, title string) ([]interface{}, error) {
	return []interface{}{}, nil
}
