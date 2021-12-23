package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func matrimonioRoutes(r *gin.RouterGroup) {
	admin := r.Group("matrimonios")
	admin.GET("", controllers.GetMatrimonios)
	admin.POST("", controllers.CreateMatrimonio)
	admin.PUT("/:id", controllers.UpdateMatrimonio)
	admin.GET("/:id", controllers.GetMatrimonioPorID)
	admin.DELETE("/:id", controllers.DeleteMatrimonio)
}
