package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func permisoRoutes(r *gin.RouterGroup) {
	admin := r.Group("permisos")
	admin.GET("", controllers.GetPermisos)
	admin.POST("", controllers.CreatePermiso)
	admin.PUT("/:id", controllers.UpdatePermiso)
	admin.GET("/:id", controllers.GetPermisoPorId)
	admin.DELETE("/:id", controllers.DeletePermiso)
}
