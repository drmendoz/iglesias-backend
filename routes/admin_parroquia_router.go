package routes

import (
	"github.com/drmendoz/iglesias-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func administradorParroquiaRouter(r *gin.RouterGroup) {
	admin := r.Group("admin-parroquia")
	admin.Use(middlewares.RolParroquiaAdminMiddleware())
	admin.Use(middlewares.AuthMiddleWare())
	admin.Use(middlewares.ParsingTokenAdminParroquia())
	administradorParroquiaRoutes(admin)
	administradorGaritaRoutes(admin)
	parroquiaRoutes(admin)
	publicacionRoutes(admin)
	residenteRoutes(admin)
	administrativosRoutes(admin)
	itemMarketRoutes(admin)
	respuestaRoutes(admin)
	permisoRoutes(admin)
	categoriaMarketRoutes(admin)
	donacionRoutes(admin)
}
