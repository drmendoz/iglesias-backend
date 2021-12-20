package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func publicidadRoutes(r *gin.RouterGroup) {
	admin := r.Group("publicidades")
	admin.GET("", controllers.GetPublicidads)
	admin.POST("", controllers.CreatePublicidad)
	admin.PUT("/:id", controllers.UpdatePublicidad)
	admin.GET("/:id", controllers.GetPublicidadPorId)
	admin.DELETE("/:id", controllers.DeletePublicidad)
}
