package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func reservacionAreaSocialRoutes(r *gin.RouterGroup) {
	admin := r.Group("reservacion")
	admin.GET("", controllers.GetReservacionAreaSocials)
	admin.POST("", controllers.CreateReservacionAreaSocial)
	admin.PUT("/:id", controllers.UpdateReservacionAreaSocial)
	admin.GET("/:id", controllers.GetReservacionAreaSocialPorId)
	admin.DELETE("/:id", controllers.DeleteReservacionAreaSocial)
}
