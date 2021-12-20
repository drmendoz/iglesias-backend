package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func administradorGaritaRoutes(r *gin.RouterGroup) {
	admin := r.Group("admin-garita")
	admin.GET("", controllers.GetAdministradoresGarita)
	admin.POST("", controllers.CreateAdministradorGarita)
	admin.PUT("/:id", controllers.UpdateAdministradorgarita)
	admin.GET("/:id", controllers.GetAdministradorGaritaPorId)
	admin.DELETE("/:id", controllers.DeleteAdministradorGarita)
}
