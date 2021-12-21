package routes

import (
	"github.com/gin-gonic/gin"
)

func administradorMasterRouter(r *gin.RouterGroup) {
	admin := r.Group("master")
	// admin.Use(middlewares.RolMasterMiddleware())
	// admin.Use(middlewares.AuthMiddleWare())
	administradorMasterRoutes(admin)
	administradorParroquiaRoutes(admin)
	categoriaMarketRoutes(admin)
	parroquiaRoutes(admin)
	pagosMarketRoutes(admin)
	publicidadRoutes(admin)
	transaccionRoutes(admin)
	permisoRoutes(admin)
	iglesiaRoutes(admin)
}
