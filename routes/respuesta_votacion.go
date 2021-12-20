package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func respuestaVotacionRoutes(r *gin.RouterGroup) {
	admin := r.Group("respuesta-votacion")
	admin.GET("", controllers.GetRespuestaVotacions)
	admin.POST("", controllers.CreateRespuestaVotacion)
	admin.PUT("/:id", controllers.UpdateRespuestaVotacion)
	admin.GET("/:id", controllers.GetRespuestaVotacionPorId)
	admin.DELETE("/:id", controllers.DeleteRespuestaVotacion)
}
