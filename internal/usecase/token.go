package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/golang-jwt/jwt/v4"
)

type Permission string

const (
	PermissionReadUserInformation     Permission = "READ_USER_INFORMATION"
	PermissionReadBoardInformation    Permission = "READ_BOARD_INFORMATION"
	PermissionReadTreasureInformation Permission = "READ_TREASURE_INFORMATION"
	PermissionReadFavorite            Permission = "READ_FAVORITE"
	PermissionCreateBoard             Permission = "CREATE_BOARD"
	PermissionCreateArticle           Permission = "PUBLISH_POSTS"
	PermissionAppendComment           Permission = "APPEND_COMMENT"
	PermissionForwardArticleToBoard   Permission = "FORWARD_ARTICLE_TO_BOARD"
	PermissionForwardArticleToEmail   Permission = "FORWARD_ARTICLE_TO_EMAIL"
	PermissionUpdateDraft             Permission = "UPDATE_DRAFT"
	PermissionDeleteDraft             Permission = "DELETE_DRAFT"
)

func (usecase *usecase) CreateAccessTokenWithUsername(username string) string {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(usecase.globalConfig.AccessTokenExpiresAt).Unix(),
		Subject:   username,
	}

	privateKey := usecase.globalConfig.AccessTokenPrivateKey
	key, err := jwt.ParseECPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		usecase.logger.Criticalf("parse private key failed: %v", err)
		return ""
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		usecase.logger.Errorf("sign token failed: %v", err)
		return ""
	}
	return ss
}

func (usecase *usecase) GetUserIDFromToken(token string) (string, error) {
	usecase.logger.Debugf("GetUserIDFromToken token: %v", token)
	pem := usecase.globalConfig.AccessTokenPublicKey
	key, err := jwt.ParseECPublicKeyFromPEM([]byte(pem))
	if err != nil {
		usecase.logger.Criticalf("parse public key failed: %v", err)
		return "", err
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		usecase.logger.Warningf("parse token failed: %v", err)
		return "", err
	}
	if jwtToken == nil {
		usecase.logger.Warningf("jwtToken == nil")
		return "", nil
	}
	if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
		usecase.logger.Debugf("subject: %v %v", claim, jwtToken.Valid)
		return claim.Subject, nil
	}
	usecase.logger.Debugf("subject: %v", jwtToken.Valid)
	return "", fmt.Errorf("token not valid")
}

func (usecase *usecase) CheckPermission(token string, permissionID []Permission, userInfo map[string]string) error {
	for _, permission := range permissionID {
		switch permission {
		case PermissionAppendComment:
			if err := usecase.checkAppendCommentPermission(token, userInfo); err != nil {
				return err
			}
		case PermissionReadBoardInformation:
			if err := usecase.checkPermissionReadBoardSettings(token, userInfo); err != nil {
				return err
			}
		case PermissionCreateBoard:
			if err := usecase.checkCreateBoardPermission(context.Background(), token); err != nil {
				return err
			}
		case PermissionReadFavorite:
		case PermissionReadTreasureInformation:
		case PermissionReadUserInformation:
			if err := usecase.checkPermissionReadUserInformation(token, userInfo); err != nil {
				return err
			}
		case PermissionUpdateDraft:
		case PermissionDeleteDraft:
		case PermissionForwardArticleToBoard:
			if err := usecase.checkForwardArticleToBoardPermission(token, userInfo); err != nil {
				return err
			}
		case PermissionForwardArticleToEmail:
			if err := usecase.checkForwardArticleToEmailPermission(token, userInfo); err != nil {
				return err
			}
		case PermissionCreateArticle:
			if err := usecase.checkCreateArticlePermission(context.Background(), token, userInfo); err != nil {
				return err
			}
		default:
			return fmt.Errorf("undefined permission id: %s", permission)
		}
	}
	return nil
}

func (usecase *usecase) checkCreateBoardPermission(ctx context.Context, token string) error {
	userID, err := usecase.GetUserIDFromToken(token)
	if err != nil {
		return fmt.Errorf("get user id from token failed: %w", err)
	}
	user, err := usecase.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user %s failed: %w", userID, err)
	}
	if !repository.UserIsSYSOP(user) {
		return fmt.Errorf("user %s does not have SYSOP permission", userID)
	}
	return nil
}

func (usecase *usecase) checkAppendCommentPermission(token string, userInfo map[string]string) error {
	return nil
}

func (usecase *usecase) checkCreateArticlePermission(ctx context.Context, token string, userInfo map[string]string) error {
	return nil
}

func (usecase *usecase) checkForwardArticleToBoardPermission(token string, userInfo map[string]string) error {
	return nil
}

func (usecase *usecase) checkForwardArticleToEmailPermission(token string, userInfo map[string]string) error {
	return nil
}

func (usecase *usecase) checkPermissionReadUserInformation(token string, userInfo map[string]string) error {
	tokenUserID, err := usecase.GetUserIDFromToken(token)
	if err != nil {
		return fmt.Errorf("get user id from token failed: %w", err)
	}
	if tokenUserID != userInfo["user_id"] {
		return fmt.Errorf("token user id is not the same userInfo user id")
	}
	return nil
}

func (usecase *usecase) checkPermissionReadBoardSettings(token string, userInfo map[string]string) error {
	tokenUserID, err := usecase.GetUserIDFromToken(token)
	if err != nil {
		return fmt.Errorf("get user id from token failed: %w", err)
	}
	boardID, ok := userInfo["board_id"]
	if !ok || boardID == "" {
		return fmt.Errorf("board_id is required to check board read permission")
	}
	return usecase.checkBoardReadPermission(context.Background(), tokenUserID, boardID)
}

func (usecase *usecase) checkBoardReadPermission(ctx context.Context, userID, boardID string) error {
	board, err := usecase.GetBoardByID(ctx, boardID)
	if err != nil {
		return fmt.Errorf("get board %s failed: %w", boardID, err)
	}
	requiredLevel, ok := repository.BoardReadPermissionLevel(board)
	if !ok || requiredLevel == 0 {
		return nil
	}
	user, err := usecase.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user %s failed: %w", userID, err)
	}
	userLevel, ok := repository.UserPermissionLevel(user)
	if !ok {
		return fmt.Errorf("user permission level is unavailable")
	}
	if userLevel&requiredLevel == 0 {
		return fmt.Errorf("user %s has no permission to read board %s", userID, boardID)
	}
	return nil
}
