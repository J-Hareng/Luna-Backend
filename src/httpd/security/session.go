package security

import (
	"errors"
	"server/src/api/db/models"
	"server/src/helper"
	"server/src/httpd/security/caches"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func generateRAWToken(username string) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Set token expiration time
	})
	return token
}
func cryptToken(token *jwt.Token, secretKey []byte) (string, error) {
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateUserToken(UL models.UserLink, GrupID string) (string, error) {
	key := []byte(helper.GetEnvVar("SECRETKEY"))

	rawT := generateRAWToken(UL.NAME)

	T, err := cryptToken(rawT, key)
	if err != nil {
		return "", err
	}

	if ok := caches.USERCACHE.Set(T, UL, GrupID); ok {
		return T, nil
	}

	//!NEEDS TO BE REMOVED JUST FOR DEBUG
	return "", errors.New("could not insert token")
}
