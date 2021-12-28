package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/drmendoz/iglesias-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func administradorParroquiaRouter(r *gin.RouterGroup) {
	admin := r.Group("admin-parroquia")
	admin.Use(middlewares.RolParroquiaAdminMiddleware())
	admin.Use(middlewares.AuthMiddleWare())
	admin.Use(middlewares.ParsingTokenAdminParroquia())
	admin.POST("archivos", controllers.SubirArchivos)
	transaccionRoutes(admin)
	administradorParroquiaRoutes(admin)
	administradorGaritaRoutes(admin)
	parroquiaRoutes(admin)
	publicacionRoutes(admin)
	administrativosRoutes(admin)
	itemMarketRoutes(admin)
	respuestaRoutes(admin)
	permisoRoutes(admin)
	categoriaMarketRoutes(admin)
	categoriaDonacionRoutes(admin)
	donacionRoutes(admin)
	fielRoutes(admin)
	padreRoutes(admin)
	misaRoutes(admin)
	cursoRoutes(admin)
	actividadRoutes(admin)
	horarioRoutes(admin)
	musicaRoutes(admin)
	matrimonioRoutes(admin)
	galeriaRoutes(admin)
	intencionRoutes(admin)
}
