package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func categoriaDonacionRoutes(r *gin.RouterGroup) {
	admin := r.Group("categorias-donacion")
	admin.GET("", controllers.GetCategoriaDonacions)
	admin.POST("", controllers.CreateCategoriaDonacion)
	admin.PUT("/:id", controllers.UpdateCategoriaDonacion)
	admin.GET("/:id", controllers.GetCategoriaDonacionPorId)
	admin.DELETE("/:id", controllers.DeleteCategoriaDonacion)
}
