package db

import (
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// Inject injects database to gin context
func Inject(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
