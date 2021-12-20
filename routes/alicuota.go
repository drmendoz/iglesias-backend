package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func alicuotaRoutes(r *gin.RouterGroup) {
	admin := r.Group("alicuota")
	admin.GET("", controllers.GetAlicuotas)
	admin.POST("", controllers.CreateAlicuota)
	admin.POST("bulk", controllers.CreateAlicuotaBulk)
	admin.PUT("/:id", controllers.UpdateAlicuota)
	admin.GET("/:id", controllers.GetAlicuotaPorId)
	admin.DELETE("/:id", controllers.DeleteAlicuota)
}
