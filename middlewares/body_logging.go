package middlewares

import (
	"bytes"
	"io/ioutil"

	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
)

func LoggingBodyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			body, _ := ioutil.ReadAll(c.Request.Body)
			utils.Log.Info(string(body))
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
		c.Next()
	}
}
