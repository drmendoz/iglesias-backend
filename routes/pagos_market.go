package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func pagosMarketRoutes(r *gin.RouterGroup) {
	admin := r.Group("suscripciones")
	admin.GET("", controllers.GetPagoMarkets)
	admin.POST("", controllers.CreatePagoMarket)
	admin.PUT("/:id", controllers.UpdatePagoMarket)
	admin.GET("/:id", controllers.GetPagoMarketPorId)
	admin.DELETE("/:id", controllers.DeletePagoMarket)
}
