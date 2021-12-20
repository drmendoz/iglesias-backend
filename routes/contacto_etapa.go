package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func contactoEtapaRoutes(r *gin.RouterGroup) {
	admin := r.Group("contactos")
	admin.GET("", controllers.GetContactoEtapas)
	admin.POST("", controllers.CreateContactoEtapa)
	admin.PUT("/:id", controllers.UpdateContactoEtapa)
	admin.GET("/:id", controllers.GetContactoEtapaPorId)
	admin.DELETE("/:id", controllers.DeleteContactoEtapa)
}
