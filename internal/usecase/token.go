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

func (t *usecase) GetUserIdFromToken(token string) (string, error) {
	t.logger.Debugf("getUserIdFromToken token: %v", token)
	pem := t.globalConfig.AccessTokenPublicKey
	key, err := jwt.ParseECPublicKeyFromPEM([]byte(pem))
	if err != nil {
		t.logger.Criticalf("parse public key failed: %v", err)
		return "", err
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{},
		func(token *jwt.Token) (i interface{}, e error) {
			return key, nil
		})
	if err != nil {
		t.logger.Warningf("parse token failed: %v", err)
		return "", err
	}

	if jwtToken == nil {
		t.logger.Warningf("jwtToken == nil")
		return "", nil
	}

	// logger.Debugf("getUserIdFromToken jwtToken: %v %v", jwtToken, err)
	if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
		t.logger.Debugf("subject: %v %v", claim, jwtToken.Valid)
		return claim.Subject, nil
		// return "", nil
	}
	t.logger.Debugf("subject: %v", jwtToken.Valid)
	return "", fmt.Errorf("token not valid")
}

func (usecase *usecase) CheckPermission(token string, permissionId []Permission, userInfo map[string]string) error {
	return nil
}
