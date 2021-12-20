package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func visualizacionRoutes(r *gin.RouterGroup) {
	admin := r.Group("visualizaciones")
	admin.POST(":modulo", controllers.ActualizarVisualizacion)
	admin.GET("", controllers.GetNotificacionesRequest)
}
