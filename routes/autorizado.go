package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func autorizadoRoutes(r *gin.RouterGroup) {
	admin := r.Group("autorizados")
	admin.GET("", controllers.GetAutorizados)
	admin.POST("", controllers.CreateAutorizado)
	admin.PUT("/:id", controllers.UpdateAutorizado)
	admin.GET("/:id", controllers.GetAutorizadoPorId)
	admin.DELETE("/:id", controllers.DeleteAutorizado)
}
