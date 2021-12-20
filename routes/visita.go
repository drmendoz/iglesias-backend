package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func visitaRoutes(r *gin.RouterGroup) {
	admin := r.Group("visita")
	admin.GET("", controllers.GetVisitas)
	admin.POST("", controllers.CreateVisita)
	admin.PUT("/:id", controllers.UpdateVisita)
	admin.GET("/:id", controllers.GetVisitaPorId)
	admin.DELETE("/:id", controllers.DeleteVisita)
	admin.POST("/:id/contestar", controllers.ContestarVisita)
}
