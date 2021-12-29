package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func administradorMasterRouter(r *gin.RouterGroup) {
	admin := r.Group("master")
	// admin.Use(middlewares.RolMasterMiddleware())
	// admin.Use(middlewares.AuthMiddleWare())

	admin.POST("archivos", controllers.SubirArchivos)
	administradorMasterRoutes(admin)
	administradorParroquiaRoutes(admin)
	categoriaMarketRoutes(admin)
	categoriaDonacionRoutes(admin)
	parroquiaRoutes(admin)
	pagosMarketRoutes(admin)
	publicidadRoutes(admin)
	transaccionRoutes(admin)
	permisoRoutes(admin)
	iglesiaRoutes(admin)
	donacionRoutes(admin)
	admin.GET("fieles", controllers.GetFieles)
	admin.PUT("/parroquia/modulos/:id", controllers.UpdateModulosParroquia)
}
