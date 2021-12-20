package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func respuestaRoutes(r *gin.RouterGroup) {
	admin := r.Group("respuesta")
	admin.GET("", controllers.GetRespuestas)
	admin.POST("", controllers.CreateRespuesta)
	admin.PUT("/:id", controllers.UpdateRespuesta)
	admin.GET("/:id", controllers.GetRespuestaPorId)
	admin.DELETE("/:id", controllers.DeleteRespuesta)
}
