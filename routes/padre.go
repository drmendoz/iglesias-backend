package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func padreRoutes(r *gin.RouterGroup) {
	admin := r.Group("padres")
	admin.GET("", controllers.GetPadres)
	admin.POST("", controllers.CreatePadre)
	admin.PUT("/:id", controllers.UpdatePadre)
	admin.GET("/:id", controllers.GetPadrePorId)
	admin.DELETE("/:id", controllers.DeletePadre)
}
