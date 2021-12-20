package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func administradorMasterRoutes(r *gin.RouterGroup) {
	admin := r.Group("master")
	//admin.Use(middlewares.RolMasterMiddleware())
	//admin.Use(middlewares.AuthMiddleWare())
	admin.GET("", controllers.GetAdministradores)
	admin.POST("", controllers.CreateAdministrador)
	admin.PUT("/:id", controllers.UpdateAdministrador)
	admin.GET("/:id", controllers.GetAdministradorPorId)
	admin.DELETE("/:id", controllers.DeleteAdministrador)
}
