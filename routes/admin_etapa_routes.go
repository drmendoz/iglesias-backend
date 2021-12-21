package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func administradorEtapaRoutes(r *gin.RouterGroup) {
	admin := r.Group("admin-etapa")
	admin.GET("", controllers.GetAdministradoresEtapa)
	admin.POST("", controllers.CreateAdministradorEtapa)
	admin.PUT("/:id", controllers.UpdateAdministradorEtapa)
	admin.GET("/:id", controllers.GetAdministradorEtapaPorId)
	admin.DELETE("/:id", controllers.DeleteAdministradorEtapa)
}
