package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func fielRoutes(r *gin.RouterGroup) {
	admin := r.Group("fiel")
	admin.GET("", controllers.GetFieles)
	admin.POST("", controllers.CreateFiel)
	admin.PUT("/:id", controllers.UpdateFiel)
	admin.GET("/:id", controllers.GetFielPorId)
	admin.DELETE("/:id", controllers.DeleteFiel)
}
