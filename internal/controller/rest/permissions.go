package rest

import (
	"github.com/dgrijalva/jwt-go"

	"fmt"
	"time"
)

type permission string

const (
	PermissionReadUserInformation     permission = "READ_USER_INFORMATION"
	PermissionReadBoardInformation    permission = "READ_BOARD_INFORMATION"
	PermissionReadTreasureInformation permission = "READ_TREASURE_INFORMATION"
	PermissionReadFavorite            permission = "READ_FAVORITE"
)

func checkTokenPermission(token string, permissionId []permission, userInfo map[string]string) error {
	return nil
}

func (rest *restHandler) getAccessTokenWithUsername(username string) string {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(rest.globalConfig.AccessTokenExpiresAt).Unix(),
		// Issuer:    "test",
		Subject: username,
	}

	// TODO: Setting me in config
	// openssl ecparam -name prime256v1 -genkey -noout -out pkey
	privateKey := rest.globalConfig.AccessTokenPrivateKey

	key, err := jwt.ParseECPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		rest.logger.Criticalf("parse private key failed: %v", err)
		return ""
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		rest.logger.Errorf("sign token failed: %v", err)
		return ""
	}
	return ss
}

func (rest *restHandler) getUserIdFromToken(token string) (string, error) {
	rest.logger.Debugf("getUserIdFromToken token: %v", token)
	pem := rest.globalConfig.AccessTokenPublicKey
	key, err := jwt.ParseECPublicKeyFromPEM([]byte(pem))
	if err != nil {
		rest.logger.Criticalf("parse public key failed: %v", err)
		return "", err
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{},
		func(token *jwt.Token) (i interface{}, e error) {
			return key, nil
		})
	if err != nil {
		rest.logger.Warningf("parse token failed: %v", err)
		return "", err
	}

	if jwtToken == nil {
		rest.logger.Warningf("jwtToken == nil")
		return "", nil
	}

	// logger.Debugf("getUserIdFromToken jwtToken: %v %v", jwtToken, err)
	if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
		rest.logger.Debugf("subject: %v %v", claim, jwtToken.Valid)
		return claim.Subject, nil
		// return "", nil
	}
	rest.logger.Debugf("subject: %v", jwtToken.Valid)
	return "", fmt.Errorf("token not valid")

}
