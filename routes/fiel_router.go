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
	historiaRoutes(res)
	tarjetaRoutes(res)
	respuestaRoutes(res)
}
