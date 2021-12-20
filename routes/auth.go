package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func authRoutes(r *gin.RouterGroup) {
	auth := r.Group("auth")
	auth.POST("/admin-master", controllers.LoginAdminMaster)
	auth.POST("/admin-garita", controllers.LoginAdminGarita)
	auth.POST("/admin-etapa", controllers.LoginAdminParroquia)
	auth.POST("/residente", controllers.LoginResidente)
	auth.POST("/residente/cambio", controllers.CambioDeContrasenaResidente)
	auth.POST("/recover/:rol", controllers.EnviarCodigoTemporal)
	auth.POST("/recover/:rol/cambio", controllers.CambioDeContrasena)
}
