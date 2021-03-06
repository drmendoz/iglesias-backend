package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/drmendoz/iglesias-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func FielRouter(r *gin.RouterGroup) {
	res := r.Group("fiel")
	res.Use(middlewares.RolFielMiddleware())
	res.Use(middlewares.AuthMiddleWare())
	res.Use(middlewares.ParsingTokenFiel())
	publicacionRoutes(res)
	res.GET("/categorias-donacion", controllers.GetCategoriaDonacions)
	res.POST("/notificacion", controllers.UpdateTokenNotificacion)
	res.GET("/administrativo", controllers.GetAdministrativos)
	res.GET("/publicidades", controllers.GetPublicidads)
	res.POST("/cambiarContrasena", controllers.CambiarContrasenaFiel)
	res.POST("/perfil", controllers.EditarImagenPerfilFiel)
	res.GET("/perfil", controllers.GetInformacionPerfil)
	res.POST("/imagen", controllers.CrearArchivo)
	res.POST("/imagenes", controllers.CrearArchivos)
	res.GET("/categorias", controllers.GetCategoriaMarkets)
	res.GET("/suscripciones", controllers.VerificarSuscripcionFiel)
	res.POST("/suscripciones", controllers.CrearSuscripcion)
	res.DELETE("/suscripciones", controllers.AnularSuscripcion)
	res.POST("/emprendimientos", controllers.CreateEmprendimiento)
	res.GET("/emprendimientos", controllers.ObtenerEmprendimientos)
	res.GET("/mis-emprendimientos", controllers.ObtenerEmprendimientosUsuarios)
	res.GET("/fieles/:id/emprendimientos", controllers.ObtenerEmprendimientoFiel)
	res.GET("/emprendimientos/:id", controllers.ObtenerEmprendimientosPorId)
	res.PUT("/emprendimientos/:id", controllers.ActualizarEmprendimiento)
	res.DELETE("/emprendimientos/:id", controllers.DeleteEmprendimiento)
	res.POST("/cerrar-sesion", controllers.CerrarSesion)
	res.POST("/causas-beneficas/:id/donar", controllers.AportarDonacion)
	res.GET("/causas-beneficas/:id", controllers.GetDonacionPorID)
	res.GET("/cursos", controllers.GetCursos)
	res.GET("/cursos/:id", controllers.GetCursoPorID)
	res.POST("/cursos/:id/inscribir", controllers.InscribirCurso)
	res.GET("/misas", controllers.GetMisas)
	res.GET("/actividades", controllers.GetActividads)
	res.GET("/actividades/:id", controllers.GetActividadPorId)
	res.GET("/horarios", controllers.GetHorarios)
	res.GET("/horarios/:id", controllers.GetHorarioPorId)
	res.GET("/musicas-usuario", controllers.GetMusicaDeUsuario)
	res.GET("/matrimonios-usuario", controllers.GetMatrimoniosUsuario)
	historiaRoutes(res)
	matrimonioRoutes(res)
	musicaRoutes(res)
	tarjetaRoutes(res)
	respuestaRoutes(res)
	intencionRoutes(res)
}
