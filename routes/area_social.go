package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func areaSocialRoutes(r *gin.RouterGroup) {
	admin := r.Group("area-social")
	admin.GET("", controllers.GetAreaSocials)
	admin.POST("", controllers.CreateAreaSocial)
	admin.PUT("/:id", controllers.UpdateAreaSocial)
	admin.GET("/:id", controllers.GetAreaSocialPorId)
	admin.DELETE("/:id", controllers.DeleteAreaSocial)
	admin.GET("/:id/horarios", controllers.GetHorariosPorAreaSocial)
	admin.GET("/:id/horarios/disponibles", controllers.GetHorarioDisponiblesAreaSocial)
}
