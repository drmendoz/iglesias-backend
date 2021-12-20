package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func buzonRoutes(r *gin.RouterGroup) {
	routes := r.Group("buzon")
	routes.GET("", controllers.GetEntradasBuzon)
	routes.GET("/:id_buzon", controllers.GetMensajesBuzones)
	routes.POST("", controllers.CreateBuzon)
	routes.POST("/:id_buzon/responder", controllers.CreateBuzon)
	routes.POST("/:id_buzon/responder-privado", controllers.ResponderRespuestaBuzonPrivado)
	routes.POST("/:id_buzon/archivos", controllers.CreateArchivosBuzon)
}
