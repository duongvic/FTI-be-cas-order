package logging

import (
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// Inject injects logger to gin context
func Inject(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", logger)
		c.Next()
	}
}
