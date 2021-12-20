package middlewares

import (
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
)

func LoggingErrorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			utils.Log.Warn(err)
		}
	}
}
