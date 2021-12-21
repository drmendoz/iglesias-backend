package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func parroquiaRoutes(r *gin.RouterGroup) {
	admin := r.Group("parroquias")
	admin.GET("", controllers.GetParroquias)
	admin.POST("", controllers.CreateParroquia)
	admin.PUT("/:id", controllers.UpdateParroquia)
	admin.GET("/:id", controllers.GetParroquiaPorId)
	admin.DELETE("/:id", controllers.DeleteParroquia)
}
