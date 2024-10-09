package security

import (
	"errors"
	"fmt"
	"net/http"
	"server/src/helper"
	"server/src/httpd/security/caches"

	"github.com/gin-gonic/gin"
)

const C_USER_KEY string = "UserInThisSession"

func ConditionToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := helper.GetVerfToken(c)
		if caches.USERCACHE.IsNil() {
			fmt.Println("UserCache is nil")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Token available in cache"})
			c.Redirect(http.StatusFound, "http://localhost:8080/")
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Token available"})
			c.Redirect(http.StatusFound, "http://localhost:8080/")
			c.Abort()
			return
		}
		if u, ok := caches.USERCACHE.Get(token); ok { // * Token finden
			fmt.Printf("User %v loggedin with token: %v\n", u.UserLink.NAME, token)
			c.Set(C_USER_KEY, u)
			c.Next()
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Invalid"})
			c.Redirect(http.StatusFound, "http://localhost:8080/")
			c.Abort()
			return
		}
	}
}
func DecodeUserFrom_C(c *gin.Context) (caches.CacheUserData, bool) {
	u_raw, ok := c.Get(C_USER_KEY)
	if !ok {
		return caches.CacheUserData{}, false
	}
	u, ok := u_raw.(caches.CacheUserData)
	if !ok {
		return caches.CacheUserData{}, false
	}
	return u, true
}
func ShouldDecodeFrom_C(c *gin.Context, v *caches.CacheUserData) error {
	u_raw, ok := c.Get(C_USER_KEY)
	if !ok {

		return errors.New("user not found")
	}
	u, ok := u_raw.(caches.CacheUserData)
	if !ok {
		return errors.New("user not found")
	}
	*v = u
	return nil
}
