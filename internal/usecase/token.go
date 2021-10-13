package usecase

import (
	"context"
	// "errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Permission string

const (
	PermissionReadUserInformation     Permission = "READ_USER_INFORMATION"
	PermissionReadBoardInformation    Permission = "READ_BOARD_INFORMATION"
	PermissionReadTreasureInformation Permission = "READ_TREASURE_INFORMATION"
	PermissionReadFavorite            Permission = "READ_FAVORITE"
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
		// Issuer:    "test",
		Subject: username,
	}

	// TODO: Setting me in config
	// openssl ecparam -name prime256v1 -genkey -noout -out pkey
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

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{},
		func(token *jwt.Token) (i interface{}, e error) {
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

	// logger.Debugf("GetUserIDFromToken jwtToken: %v %v", jwtToken, err)
	if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
		usecase.logger.Debugf("subject: %v %v", claim, jwtToken.Valid)
		return claim.Subject, nil
		// return "", nil
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

func (usecase *usecase) checkAppendCommentPermission(token string, userInfo map[string]string) error {
	// boardID := userInfo["board_id"]
	// userID := userInfo["user_id"]

	// TODO: 判斷在該版是否被水桶
	// TODO: 判斷該版是否允許推文
	// TODO: 判斷該文章是否鎖文

	return nil
}

func (usecase *usecase) checkCreateArticlePermission(ctx context.Context, token string, userInfo map[string]string) error {
	// TODO:
	// get board data and check whether board allow create articles
	// boardID, ok := userInfo["board_id"]
	// if !ok {
	// 	return errors.New("no board_id key") // todo: define error
	// }
	// boardLimit, err := usecase.repo.GetBoardPostsLimit(ctx, boardID)
	// if err != nil {
	// 	return fmt.Errorf("create article permission failed: %w", err)
	// }
	// if !boardLimit.EnableNewPost() {
	// 	return errors.New("this board not allow create new article")
	// }

	// get board ban list and check whether user on the list
	// TODO: repo 新增各板水桶名單
	// get global ban list and check whether user on the list
	// TODO: repo 新增全站水桶名單

	return nil
}

// This function checks the user has permission that can forward the target article to another board.
func (usecase *usecase) checkForwardArticleToBoardPermission(token string, userInfo map[string]string) error {
	// boardID := userInfo["board_id"]
	// toBoard := userInfo["to_board"]
	// userID := userInfo["user_id"]

	// TODO: 判斷是否有轉錄的權限
	// TODO: 判斷在該版是否允許發文
	// TODO: 判斷轉錄的版是否允許發文
	// TODO: 判斷在該版是否被水桶
	// TODO: 判斷轉錄的次數上限
	// TODO: 判斷轉錄的版跟現在的版是不同的
	// TODO: 判斷冷卻時間
	// TODO: 判斷 CAPTCHA 驗證是否通過

	return nil
}

// This function checks the user has permission that can forward the target article to the target email address.
// Implementation should note that the target email is not a private email or an unresolved address.
func (usecase *usecase) checkForwardArticleToEmailPermission(token string, userInfo map[string]string) error {
	// boardID := userInfo["board_id"]
	// toEmail := userInfo["to_email"]
	// userID := userInfo["user_id"]

	// TODO: 確認Email是可以轉發的
	// TODO: 判斷是否有轉錄的權限
	// TODO: 判斷在該版是否允許發文
	// TODO: 判斷轉錄的版是否允許發文
	// TODO: 判斷在該版是否被水桶
	// TODO: 判斷轉錄的次數上限
	// TODO: 判斷轉錄的版跟現在的版是不同的
	// TODO: 判斷冷卻時間
	// TODO: 判斷 CAPTCHA 驗證是否通過

	return nil
}

func (usecase *usecase) checkPermissionReadUserInformation(token string, userInfo map[string]string) error {
	// TODO: 判斷管理群的權限
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
	// TODO: 判斷管理群的權限
	_, err := usecase.GetUserIDFromToken(token)
	if err != nil {
		return fmt.Errorf("get user id from token failed: %w", err)
	}

	return nil
}
