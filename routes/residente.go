package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func residenteRoutes(r *gin.RouterGroup) {
	admin := r.Group("residente")
	//admin.Use(middlewares.RolMasterMiddleware())
	//admin.Use(middlewares.AuthMiddleWare())
	admin.GET("", controllers.GetFiel)
	admin.POST("", controllers.CreateFiel)
	admin.PUT("/:id", controllers.UpdateFiel)
	admin.GET("/:id", controllers.GetFielPorId)
	admin.DELETE("/:id", controllers.DeleteFiel)
}
