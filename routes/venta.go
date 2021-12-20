package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func ventaRoutes(r *gin.RouterGroup) {
	admin := r.Group("ventas")
	admin.GET("", controllers.GetAllVentas)
	admin.POST("", controllers.CreateVenta)
	admin.GET("/:id", controllers.GetVentaById)
	admin.DELETE("/:id", controllers.DeleteVenta)
}
