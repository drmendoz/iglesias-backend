package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func actividadRoutes(r *gin.RouterGroup) {
	admin := r.Group("actividades")
	admin.GET("", controllers.GetActividads)
	admin.POST("", controllers.CreateActividad)
	admin.PUT("/:id", controllers.UpdateActividad)
	admin.GET("/:id", controllers.GetActividadPorId)
	admin.DELETE("/:id", controllers.DeleteActividad)
}
