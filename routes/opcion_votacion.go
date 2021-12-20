package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func opcionVotacionRoutes(r *gin.RouterGroup) {
	admin := r.Group("opcion-votacion")
	admin.GET("", controllers.GetOpcionVotacions)
	admin.POST("", controllers.CreateOpcionVotacion)
	admin.PUT("/:id", controllers.UpdateOpcionVotacion)
	admin.GET("/:id", controllers.GetVotacionPorId)
	admin.DELETE("/:id", controllers.DeleteOpcionVotacion)
}
