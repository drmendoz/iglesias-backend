package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func transaccionRoutes(r *gin.RouterGroup) {
	admin := r.Group("transacciones")
	//admin.Use(middlewares.RolMasterMiddleware())
	//admin.Use(middlewares.AuthMiddleWare())
	admin.GET("", controllers.GetTransaccion)
	admin.GET(":id", controllers.GetTransaccionPorId)
	admin.POST(":id/devolver", controllers.DevolverTransaccion)
}
