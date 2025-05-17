package middleware

import (
	"7-solutions/repository"
	"7-solutions/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func JWTAuth(secret string, userRepo repository.UsersRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenStr, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
			c.Abort()
			return
		}

		user, err := userRepo.GetByID(c.Request.Context(), userID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user does not exist (possibly deleted)"})
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking user"})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID.Hex())
		c.Next()
	}
}
