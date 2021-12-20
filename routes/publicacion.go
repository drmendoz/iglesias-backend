package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func publicacionRoutes(r *gin.RouterGroup) {
	admin := r.Group("noticias")
	admin.GET("", controllers.GetPublicacions)
	admin.POST("", controllers.CreatePublicacion)
	admin.PUT("/:id", controllers.UpdatePublicacion)
	admin.GET("/:id", controllers.GetPublicacionPorId)
	admin.DELETE("/:id", controllers.DeletePublicacion)
	admin.POST("/:id/leido", controllers.MarcarPublicacionLeida)
}
