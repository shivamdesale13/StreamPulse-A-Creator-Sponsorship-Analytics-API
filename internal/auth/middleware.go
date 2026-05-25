package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ctxUserID = "user_id"
	ctxEmail  = "user_email"
)

func Middleware(jwt *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwt.Verify(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxEmail, claims.Email)
		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	v, _ := c.Get(ctxUserID)
	id, _ := v.(string)
	return id
}
