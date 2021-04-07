package usecase

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Permission string

const (
	PermissionReadUserInformation     Permission = "READ_USER_INFORMATION"
	PermissionReadBoardInformation    Permission = "READ_BOARD_INFORMATION"
	PermissionReadTreasureInformation Permission = "READ_TREASURE_INFORMATION"
	PermissionReadFavorite            Permission = "READ_FAVORITE"
	PermissionAppendComment           Permission = "APPEND_COMMENT"
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
		case PermissionReadFavorite:
		case PermissionReadTreasureInformation:
		case PermissionReadUserInformation:
		default:
			return fmt.Errorf("undefined permission id:%s", permission)
		}
	}

	return nil
}

func (usecase *usecase) checkAppendCommentPermission(token string, userInfo map[string]string) error {
	//boardID := userInfo["board_id"]
	//userID := userInfo["user_id"]

	// TODO: 判斷在該版是否被水桶
	// TODO: 判斷該版是否允許推文
	// TODO: 判斷該文章是否鎖文

	return nil
}
