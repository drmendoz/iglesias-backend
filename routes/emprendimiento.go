package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func itemMarketRoutes(r *gin.RouterGroup) {
	admin := r.Group("emprendimiento")
	admin.GET("", controllers.ObtenerEmprendimientos)
	admin.POST("", controllers.CreateEmprendimiento)
	admin.PUT("/:id", controllers.ActualizarEmprendimiento)
	admin.GET("/:id", controllers.ObtenerEmprendimientosPorId)
	admin.DELETE("/:id", controllers.DeleteEmprendimiento)
}
