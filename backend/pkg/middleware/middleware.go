package middleware

import (
	"net/http"

	"github.com/ChayanDass/beneficiary-manager/pkg/db"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CORSMiddleware is a middleware function for CORS.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("X-Username")
		password := c.GetHeader("X-Password")

		if username == "" || password == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing username or password in headers",
			})
			return
		}

		// Find the user by username
		var user models.User
		if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid username",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		// Set in context
		c.Set("username", username)
		c.Set("password", password)
		c.Set("user_id", user.ID)
		c.Next()
	}
}
