package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func donacionRoutes(r *gin.RouterGroup) {
	admin := r.Group("causas-beneficas")
	admin.GET("", controllers.GetDonacions)
	admin.GET("/:id/aportaciones", controllers.GetAportacionesDeDonacion)
	admin.POST("", controllers.CreateDonacion)
	admin.PUT("/:id", controllers.UpdateDonacion)
	admin.GET("/:id", controllers.GetDonacionPorID)
	admin.DELETE("/:id", controllers.DeleteDonacion)
}
