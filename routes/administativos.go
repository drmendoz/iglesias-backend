package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func administrativosRoutes(r *gin.RouterGroup) {
	admin := r.Group("administrativos")
	admin.GET("", controllers.GetAdministrativos)
	admin.POST("", controllers.CreateAdministrativo)
	admin.PUT("/:id", controllers.UpdateAdministrativo)
	admin.GET("/:id", controllers.GetAdministrativoPorId)
	admin.DELETE("/:id", controllers.DeleteAdministrativo)
}
