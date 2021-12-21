package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func etapaRoutes(r *gin.RouterGroup) {
	admin := r.Group("etapa")
	admin.GET("", controllers.GetParroquias)
	admin.POST("", controllers.CreateParroquia)
	admin.PUT("/:id", controllers.UpdateParroquia)
	admin.GET("/:id", controllers.GetParroquiaPorId)
	admin.DELETE("/:id", controllers.DeleteParroquia)
}
