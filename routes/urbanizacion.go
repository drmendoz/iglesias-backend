package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func urbanizacionRoutes(r *gin.RouterGroup) {
	admin := r.Group("urbanizacion")
	//admin.Use(middlewares.RolMasterMiddleware())
	//admin.Use(middlewares.AuthMiddleWare())
	admin.GET("", controllers.GetUrbanizaciones)
	admin.POST("", controllers.CreateUrbanizacion)
	admin.PUT("/:id", controllers.UpdateUrbanizacion)
	admin.GET("/:id", controllers.GetUrbanizacionPorId)
	admin.DELETE("/:id", controllers.DeleteUrbanizacion)
	admin.GET("/:id/casas", controllers.GetCasasPorUrbanizacion)
	admin.GET("/:id/etapas", controllers.GetEtapasPorUrbanizacion)
}
