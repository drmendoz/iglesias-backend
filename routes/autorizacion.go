package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func autorizacionRoutes(r *gin.RouterGroup) {
	admin := r.Group("autorizacion")
	admin.GET("", controllers.GetAutorizaciones)
	admin.POST("", controllers.CreateAutorizacion)
	admin.PUT("/:id", controllers.UpdateAutorizacion)
	admin.GET("/:id", controllers.GetAutorizacionPorId)
	admin.DELETE("/:id", controllers.DeleteAutorizacion)

}
