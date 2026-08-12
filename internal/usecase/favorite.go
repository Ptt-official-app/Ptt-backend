package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Ptt-official-app/go-bbs"
)

var (
	ErrFavoriteUnauthorized = errors.New("favorite modification unauthorized")
	ErrFavoriteForbidden    = errors.New("favorite modification forbidden")
	ErrFavoriteUserNotFound = errors.New("favorite user not found")
	ErrInvalidFavorite      = errors.New("invalid favorite")
)

func (usecase *usecase) AddUserFavorite(ctx context.Context, token, userID, favoriteType, boardID, title string) ([]interface{}, error) {
	tokenUserID, err := usecase.GetUserIDFromToken(token)
	if err != nil || tokenUserID == "" {
		return nil, fmt.Errorf("%w: %v", ErrFavoriteUnauthorized, err)
	}
	if !strings.EqualFold(tokenUserID, userID) {
		return nil, fmt.Errorf("%w: token user %s cannot modify favorites for %s", ErrFavoriteForbidden, tokenUserID, userID)
	}
	if _, err := usecase.GetUserByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFavoriteUserNotFound, err)
	}

	options := bbs.FavoriteCreateOptions{BoardID: boardID, Title: title}
	switch strings.ToLower(favoriteType) {
	case "line":
		options.Type = bbs.FavoriteTypeLine
	case "folder":
		if title == "" {
			return nil, fmt.Errorf("%w: folder title is required", ErrInvalidFavorite)
		}
		options.Type = bbs.FavoriteTypeFolder
	case "board":
		if boardID == "" {
			return nil, fmt.Errorf("%w: board_id is required", ErrInvalidFavorite)
		}
		if err := usecase.CheckPermission(token, []Permission{PermissionReadBoardInformation}, map[string]string{"board_id": boardID}); err != nil {
			return nil, fmt.Errorf("%w: cannot favorite board %s: %v", ErrFavoriteForbidden, boardID, err)
		}
		options.Type = bbs.FavoriteTypeBoard
	default:
		return nil, fmt.Errorf("%w: unsupported type %q", ErrInvalidFavorite, favoriteType)
	}

	if _, err := usecase.repo.AddUserFavorite(ctx, userID, options); err != nil {
		return nil, fmt.Errorf("add favorite: %w", err)
	}
	return usecase.GetUserFavorites(ctx, userID)
}
