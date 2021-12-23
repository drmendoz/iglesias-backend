package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func intencionRoutes(r *gin.RouterGroup) {
	admin := r.Group("intenciones")
	admin.GET("", controllers.GetIntenciones)
	admin.POST("", controllers.CreateIntencion)
	admin.PUT("/:id", controllers.UpdateIntencion)
	admin.GET("/:id", controllers.GetIntencionPorId)
	admin.DELETE("/:id", controllers.DeleteIntencion)
}
