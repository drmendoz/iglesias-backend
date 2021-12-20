package routes

import (
	"github.com/drmendoz/iglesias-backend/controllers"
	"github.com/drmendoz/iglesias-backend/middlewares"
	"github.com/drmendoz/iglesias-backend/sockets"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var R *gin.Engine

func init() {
	setLogging()
	if utils.Viper.GetBool("PROD") {
		gin.SetMode(gin.ReleaseMode)
	}

	R = gin.Default()
	//R.Use(middlewares.LoggingBodyMiddleware())
	R.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           86400,
	}))
	router := R.Group("/api/v1")
	router.Static("/public", "./public")
	router.POST("/alicuotas", controllers.CreateAlicuotaAutomaticamente)
	router.POST("/suscripciones", controllers.RenovarSuscripciones)
	authRoutes(router)
	administradorMasterRouter(router)
	administradorGaritaRouter(router)
	administradorEtapaRouter(router)
	ResidenteRouter(router)
	DataRouter(router)
	R.GET("/socket.io/*any", gin.WrapH(sockets.ServerVisita))
	R.POST("/socket.io/*any", gin.WrapH(sockets.ServerVisita))
	R.Use(middlewares.LoggingErrorsMiddleware())
}
