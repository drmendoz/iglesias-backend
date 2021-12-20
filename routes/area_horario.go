package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func areaHorarioRoutes(r *gin.RouterGroup) {
	admin := r.Group("area-horarios")
	admin.GET("", controllers.GetAreaHorarios)
	admin.POST("", controllers.CreateAreaHorario)
	admin.PUT("/:id", controllers.UpdateAreaHorario)
	admin.GET("/:id", controllers.GetAreaHorarioPorId)
	admin.DELETE("/:id", controllers.DeleteAreaHorario)
}
