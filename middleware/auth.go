// middleware/auth.go

package middleware

import (
	"fmt"
	"library/task/models"
	"library/task/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("user")
		fmt.Println("token print", tokenString)
		if tokenString == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "token is missing"})
			c.Abort()
			return
		}

		token, err := utils.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*utils.Claims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		user, err := models.GetUserByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func Authorize(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)

		fmt.Println(user)

		if user.Role != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to proceed further"})
			c.Abort()
			return
		}

		c.Next()
	}
}
