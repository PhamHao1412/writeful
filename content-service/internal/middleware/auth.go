package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type JWTClaims struct {
	TokenType string `json:"type"`
	UserID    string `json:"user_id"`
}

type JWTVerifier interface {
	VerifyToken(token string) (JWTClaims, error)
}

func AuthMiddleware(verifier JWTVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		claims, err := verifier.VerifyToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		if claims.TokenType != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token type"})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
