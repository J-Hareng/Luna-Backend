package helper

import (
	"net/http"
	"strings"

	"errors"

	"github.com/gin-gonic/gin"
)

func GetVerfToken(c *gin.Context) (string, error) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return "", errors.New("no value in auth header")
	}

	token := strings.Split(authorization, " ")
	if len(token) != 2 || token[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		c.Abort()
		return "", errors.New("invalid authorization header format")
	}

	return token[1], nil
}
