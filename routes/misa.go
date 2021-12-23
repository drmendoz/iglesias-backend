package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func misaRoutes(r *gin.RouterGroup) {
	admin := r.Group("misas")
	admin.GET("", controllers.GetMisas)
	admin.POST("", controllers.CreateMisa)
	admin.PUT("/:id", controllers.UpdateMisa)
	admin.GET("/:id", controllers.GetMisaPorId)
	admin.DELETE("/:id", controllers.DeleteMisa)
}
