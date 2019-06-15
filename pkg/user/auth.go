package user

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JwtAuthentication is middleware for checking jwt inside header and setting user id
func JwtAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		notAuth := []string{"/api/v1/register", "/api/v1/login"}
		requestPath := c.Request.URL.Path

		for _, value := range notAuth {
			if value == requestPath {
				c.Next()
				return
			}
		}

		tokenHeader := c.GetHeader("Authorization")

		if tokenHeader == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "Fail", "message": "Missing Authorization token"})
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "Fail", "message": "Token is malformed"})
			return
		}

		tokenPart := splitted[1]
		tk := &Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "Fail", "message": err.Error()})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "Fail", "message": "Token is not valid"})
			return
		}

		c.Set("user", tk.UserId)
		c.Next()
	}
}
