package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func authRoutes(r *gin.RouterGroup) {
	auth := r.Group("auth")
	auth.POST("/admin-master", controllers.LoginAdminMaster)
	auth.POST("/admin-parroquia", controllers.LoginAdminParroquia)
	auth.POST("/fiel", controllers.LoginFiel)
	auth.POST("/fiel/cambio", controllers.CambioDeContrasenaFiel)
	auth.POST("/recover/:rol", controllers.EnviarCodigoTemporal)
	auth.POST("/recover/:rol/cambio", controllers.CambioDeContrasena)
}
