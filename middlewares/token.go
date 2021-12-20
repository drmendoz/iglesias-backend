package middlewares

import "github.com/gin-gonic/gin"

func RolMasterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("rol", "admin-master")
		c.Next()
	}
}

func RolEtapaAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("rol", "admin-etapa")
		c.Next()
	}
}

func RolGaritaAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("rol", "admin-garita")
		c.Next()
	}
}

func RolResidenteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("rol", "residente")
		c.Next()
	}
}
