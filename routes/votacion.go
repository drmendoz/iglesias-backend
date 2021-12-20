package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func votacionRoutes(r *gin.RouterGroup) {
	admin := r.Group("votacion")
	admin.GET("", controllers.GetVotacions)
	admin.POST("", controllers.CreateVotacion)
	admin.PUT("/:id", controllers.UpdateVotacion)
	admin.GET("/:id", controllers.GetVotacionPorId)
	admin.DELETE("/:id", controllers.DeleteVotacion)
}
