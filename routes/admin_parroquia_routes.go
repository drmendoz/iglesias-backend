package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func administradorParroquiaRoutes(r *gin.RouterGroup) {
	admin := r.Group("admin-parroquias")
	admin.GET("", controllers.GetAdministradoresParroquia)
	admin.POST("", controllers.CreateAdministradorParroquia)
	admin.PUT("/:id", controllers.UpdateAdministradorParroquia)
	admin.GET("/:id", controllers.GetAdministradorParroquiaPorId)
	admin.DELETE("/:id", controllers.DeleteAdministradorParroquia)
}
