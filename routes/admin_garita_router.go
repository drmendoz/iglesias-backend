package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/drmendoz/iglesias-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func administradorGaritaRouter(r *gin.RouterGroup) {
	admin := r.Group("admin-garita")
	admin.Use(middlewares.RolGaritaAdminMiddleware())
	admin.Use(middlewares.AuthMiddleWare())
	admin.Use(middlewares.ParsingTokenAdminGarita())
	admin.GET("visita-notificaciones", controllers.NotificarVisita)
	admin.GET("visita-notificaciones/:id", controllers.NotificarVisitaPorId)
	admin.GET("/validar-autorizacion", controllers.ValidarAutorizacion)
	admin.GET("/autorizacion", controllers.GetAutorizacionesAdmin)
	admin.PUT("/autorizacion/:id", controllers.UpdateAutorizacion)
	visitaRoutes(admin)
	administradorGaritaRoutes(admin)
	areaSocialRoutes(admin)
	casaRoutes(admin)
	alicuotaRoutes(admin)
	etapaRoutes(admin)
	galeriaRoutes(admin)
	areaHorarioRoutes(admin)
	residenteRoutes(admin)
	deliveryRoutes(admin)
	autorizadoRoutes(admin)
}
