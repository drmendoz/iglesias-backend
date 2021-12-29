package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/gin-gonic/gin"
)

func cursoRoutes(r *gin.RouterGroup) {
	admin := r.Group("cursos")
	admin.GET("", controllers.GetCursos)
	admin.POST("", controllers.CreateCurso)
	admin.PUT("/:id", controllers.UpdateCurso)
	admin.GET("/:id", controllers.GetCursoPorID)
	admin.GET("/:id/inscritos", controllers.GetInscritosTotal)
	admin.DELETE("/:id", controllers.DeleteCurso)
}
