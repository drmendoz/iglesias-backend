package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func horarioRoutes(r *gin.RouterGroup) {
	admin := r.Group("horarios")
	admin.GET("", controllers.GetHorarios)
	admin.POST("", controllers.CreateHorario)
	admin.PUT("/:id", controllers.UpdateHorario)
	admin.GET("/:id", controllers.GetHorarioPorId)
	admin.DELETE("/:id", controllers.DeleteHorario)
}
