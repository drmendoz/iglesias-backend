package middlewares

import "github.com/gin-gonic/gin"

func RolMasterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("rol", "admin-master")
		c.Next()
	}
}

func RolParroquiaAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("rol", "admin-parroquia")
		c.Next()
	}
}

func RolFielMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("rol", "fiel")
		c.Next()
	}
}
