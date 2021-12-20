package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func residenteRoutes(r *gin.RouterGroup) {
	admin := r.Group("residente")
	//admin.Use(middlewares.RolMasterMiddleware())
	//admin.Use(middlewares.AuthMiddleWare())
	admin.GET("", controllers.GetResidente)
	admin.POST("", controllers.CreateResidente)
	admin.PUT("/:id", controllers.UpdateResidente)
	admin.GET("/:id", controllers.GetResidentePorId)
	admin.DELETE("/:id", controllers.DeleteResidente)
}
