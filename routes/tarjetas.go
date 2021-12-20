package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func tarjetaRoutes(r *gin.RouterGroup) {
	admin := r.Group("tarjetas")
	//admin.Use(middlewares.RolMasterMiddleware())
	//admin.Use(middlewares.AuthMiddleWare())
	admin.GET("", controllers.GetTarjetas)
	admin.DELETE(":token", controllers.DeleteTarjeta)
	admin.POST(":token/cobro", controllers.CobrarTarjeta)
}
