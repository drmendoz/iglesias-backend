package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/drmendoz/iglesias-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func administradorEtapaRouter(r *gin.RouterGroup) {
	admin := r.Group("admin-etapa")
	admin.Use(middlewares.RolEtapaAdminMiddleware())
	admin.Use(middlewares.AuthMiddleWare())
	admin.Use(middlewares.ParsingTokenAdminEtapa())
	administradorEtapaRoutes(admin)
	administradorGaritaRoutes(admin)
	urbanizacionRoutes(admin)
	etapaRoutes(admin)
	reservacionAreaSocialRoutes(admin)
	alicuotaRoutes(admin)
	areaSocialRoutes(admin)
	contactoEtapaRoutes(admin)
	galeriaRoutes(admin)
	publicacionRoutes(admin)
	admin.GET("/autorizacion", controllers.GetAutorizacionesAdmin)
	admin.GET("alicuotas", controllers.ReporteAlicuotas)
	admin.PUT("alicuotas/bulk", controllers.UpdateAlicuotaBulk)
	admin.POST("noticias-media", controllers.CreatePublicacionMedia)
	admin.GET("area_social/:id", controllers.GetAreaSocialPorID)
	residenteRoutes(admin)
	votacionRoutes(admin)
	administrativosRoutes(admin)
	camaraEtapaRoutes(admin)
	areaHorarioRoutes(admin)
	casaRoutes(admin)
	visitaRoutes(admin)
	autorizadoRoutes(admin)
	itemMarketRoutes(admin)
	mensajeRoutes(admin)
	respuestaRoutes(admin)
	permisoRoutes(admin)
	categoriaMarketRoutes(admin)
	buzonRoutes(admin)
	admin.GET("buzon-recibidos", controllers.GetBuzonesRecibidosAdminEtapa)
	admin.GET("buzon-enviados", controllers.GetBuzonesEnviados)
}
