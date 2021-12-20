package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/drmendoz/iglesias-backend/sockets"
	"github.com/gin-gonic/gin"
)

func DataRouter(r *gin.RouterGroup) {
	data := r.Group("data")
	data.GET("/casas", controllers.GetCasasCount)
	data.GET("/urbanizaciones", controllers.GetUrbanizacionesCount)
	data.GET("/residentes", controllers.GetUFielsCount)
	data.POST("/contacto", controllers.SolicitudContacto)
	data.GET("/bitacora-server", gin.WrapH(sockets.ServerVisita))
}
