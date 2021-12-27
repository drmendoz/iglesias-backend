package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func galeriaRoutes(r *gin.RouterGroup) {
	admin := r.Group("galerias")
	admin.GET("", controllers.GetGalerias)
	admin.POST("", controllers.CreateGaleria)
	admin.PUT("/:id", controllers.UpdateGaleria)
	admin.GET("/:id", controllers.GetGaleriaPorId)
	admin.DELETE("/:id", controllers.DeleteGaleria)
}
