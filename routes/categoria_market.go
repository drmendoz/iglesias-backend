package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func categoriaMarketRoutes(r *gin.RouterGroup) {
	admin := r.Group("categorias")
	//admin.Use(middlewares.RolMasterMiddleware())
	//admin.Use(middlewares.AuthMiddleWare())
	admin.GET("", controllers.GetCategoriaMarkets)
	admin.POST("", controllers.CreateCategoriaMarket)
	admin.PUT("/:id", controllers.UpdateCategoriaMarket)
	admin.GET("/:id", controllers.GetCategoriaMarketPorId)
	admin.DELETE("/:id", controllers.DeleteCategoriaMarket)
}
