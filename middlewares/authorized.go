package middlewares

import (
	"casorder/db/models"
	"casorder/utils/logging"
	"casorder/utils/mgrpc"
	"casorder/utils/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

var LOG = logging.GetLogger()

// GetAuthorizedUser Authorized blocks unauthorized requesters
func GetAuthorizedUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Invalid User"})
		c.AbortWithStatus(401)
		LOG.Error("Error: Authorization Failed")
		return nil, false
	} else {
		switch user := user.(type) {
		case *models.User:
			return user, true
		default:
			c.AbortWithStatus(403)
			LOG.Error("Error: Invalid User")
			c.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Invalid User"})
			return nil, false
		}
	}
}

func validateToken(token string) (types.JSON, error) {
	userData, err := mgrpc.VerifyToken(token)
	if err != nil {
		LOG.Error(err)
		return nil, err
	}
	return userData, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			token = c.Request.Header.Get("Authorization")
			if token == "" {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"Error": "User unauthorized"})
				c.AbortWithStatus(401)
				return
			}
		}

		tokenData, err := validateToken(token)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"Error": "Authorization failed"})
			c.AbortWithStatus(401)
			return
		}

		var user models.User
		user.Read(tokenData)

		c.Set("user", &user)
		c.Next()
	}
}
