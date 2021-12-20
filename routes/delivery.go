package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func deliveryRoutes(r *gin.RouterGroup) {
	admin := r.Group("delivery")
	admin.GET("", controllers.GetDeliverys)
	admin.POST("", controllers.CreateDelivery)
	admin.PUT("/:id", controllers.UpdateDelivery)
	admin.GET("/:id", controllers.GetDeliveryPorId)
	admin.DELETE("/:id", controllers.DeleteDelivery)
}
