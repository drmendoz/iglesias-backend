package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func mensajeRoutes(r *gin.RouterGroup) {
	admin := r.Group("mensaje")
	admin.GET("", controllers.GetMensajes)
	admin.POST("", controllers.CreateMensaje)
	admin.PUT("/:id", controllers.UpdateMensaje)
	admin.GET("/:id", controllers.GetMensajePorId)
	admin.DELETE("/:id", controllers.DeleteMensaje)
}
