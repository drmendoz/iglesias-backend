package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func casaRoutes(r *gin.RouterGroup) {
	admin := r.Group("casa")
	admin.GET("", controllers.GetCasas)
	admin.POST("", controllers.CreateCasa)
	admin.PUT("/:id", controllers.UpdateCasa)
	admin.GET("/:id", controllers.GetCasaPorId)
	admin.DELETE("/:id", controllers.DeleteCasa)
	admin.GET("/:id/residentes", controllers.GetResidentesPorCasa)
	admin.GET("/:id/alicuotas", controllers.GetAlicuotaPorCasa)
}
