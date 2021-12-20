package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func galeriaRoutes(r *gin.RouterGroup) {
	admin := r.Group("galeria")
	admin.GET("", controllers.GetImagenGalerias)
	admin.POST("", controllers.CreateImagenGaleria)
	admin.PUT("/:id", controllers.UpdateImagenGaleria)
	admin.GET("/:id", controllers.GetImagenGaleriaPorId)
	admin.DELETE("/:id", controllers.DeleteImagenGaleria)
}
