package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func iglesiaRoutes(r *gin.RouterGroup) {
	admin := r.Group("iglesias")
	admin.GET("", controllers.GetIglesiaes)
	admin.GET(":id", controllers.GetIglesiaPorId)
	admin.PUT(":id", controllers.UpdateIglesia)
	admin.DELETE(":id", controllers.DeleteIglesia)
	admin.POST("", controllers.CreateIglesia)
}
