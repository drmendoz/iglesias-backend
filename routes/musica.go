package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func musicaRoutes(r *gin.RouterGroup) {
	admin := r.Group("musicas")
	admin.GET("", controllers.GetMusicas)
	admin.POST("", controllers.CreateMusica)
	admin.PUT("/:id", controllers.UpdateMusica)
	admin.GET("/:id", controllers.GetMusicaPorID)
	admin.DELETE("/:id", controllers.DeleteMusica)
}
