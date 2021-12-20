package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func etapaRoutes(r *gin.RouterGroup) {
	admin := r.Group("etapa")
	admin.GET("", controllers.GetEtapas)
	admin.POST("", controllers.CreateEtapa)
	admin.PUT("/:id", controllers.UpdateEtapa)
	admin.GET("/:id", controllers.GetEtapaPorId)
	admin.DELETE("/:id", controllers.DeleteEtapa)
	admin.GET("/:id/casas", controllers.GetCasasPorEtapa)
	admin.GET("/:id/publicidades", controllers.GetPublicidadPorEtapaId)
}
