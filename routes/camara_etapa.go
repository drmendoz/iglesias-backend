package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func camaraEtapaRoutes(r *gin.RouterGroup) {
	admin := r.Group("camaras")
	admin.GET("", controllers.GetEtapaCamaras)
	admin.POST("", controllers.CreateEtapaCamara)
	admin.PUT("/:id", controllers.UpdateEtapaCamara)
	admin.GET("/:id", controllers.GetEtapaCamaraPorId)
	admin.DELETE("/:id", controllers.DeleteEtapaCamara)
}
