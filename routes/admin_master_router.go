package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/drmendoz/iglesias-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func administradorMasterRouter(r *gin.RouterGroup) {
	admin := r.Group("master")
	admin.Use(middlewares.RolMasterMiddleware())
	admin.Use(middlewares.AuthMiddleWare())
	administradorMasterRoutes(admin)
	administradorEtapaRoutes(admin)
	urbanizacionRoutes(admin)
	categoriaMarketRoutes(admin)
	alicuotaRoutes(admin)
	etapaRoutes(admin)
	pagosMarketRoutes(admin)
	publicidadRoutes(admin)
	transaccionRoutes(admin)
	ventaRoutes(admin)
	permisoRoutes(admin)
	admin.GET("bitacoras/etapa/:id", controllers.GetVisitasPorEtapa)
	admin.DELETE("/autorizacion/:id", controllers.DeleteAutorizacion)
}
