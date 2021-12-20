package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bradfitz/slice"
	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/tiempo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SubirHistoria(c *gin.Context) {
	idFiel := uint(c.GetInt("id_residente"))
	archivo, err := c.FormFile("historia")
	isVid := false
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Parametros de request invalidos"), nil, c, http.StatusBadRequest)
		return
	}
	nombreTmp := "public/historias/" + time.Now().Format(time.RFC3339) + archivo.Filename
	isVideo := c.Request.FormValue("isVideo")
	err = c.SaveUploadedFile(archivo, nombreTmp)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al subir historia"), nil, c, http.StatusInternalServerError)
		return
	}
	nombreFinal := ""
	if isVideo == "true" {
		isVid = true
		// 	video, err := cinema.Load(nombreTmp)
		// 	if err != nil {
		// 		_ = c.Error(err)
		// 		utils.CrearRespuesta(errors.New("Error al subir historia"), nil, c, http.StatusInternalServerError)
		// 		return
		// 	}
		// 	if video.Duration().Seconds() > 15 {
		// 		video.SetEnd(15 * time.Second)
		// 	}
		// 	nombreFinal = "public/historias/" + time.Now().Format(time.RFC3339) + ".mp4"
		// 	err = video.Render(nombreFinal)
		// 	if err != nil {
		// 		_ = c.Error(err)
		// 		utils.CrearRespuesta(errors.New("Error al subir historia"), nil, c, http.StatusInternalServerError)
		// 		return
	}
	// 	_ = os.Remove(nombreTmp)
	// } else {
	nombreFinal = nombreTmp
	//}
	fechaInicio := time.Now().In(tiempo.Local)
	fechaFinal := fechaInicio.Add(time.Hour*time.Duration(23) +
		time.Minute*time.Duration(59) +
		time.Second*time.Duration(0))
	err = models.Db.Create(&models.HistoriaFiel{Url: nombreFinal, IsVideo: isVid, FielID: idFiel, FechaFin: fechaFinal}).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al subir historia"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Historia subida exitosamente", c, http.StatusCreated)
}

type UsuarioHistoria struct {
	ID               uint          `json:"ID"`
	Key              string        `json:"key"`
	Nombre           string        `json:"name"`
	Leido            bool          `json:"hasStory"`
	Source           ImagenUsuario `json:"source"`
	FechaPublicacion time.Time     `json:"fecha_publicacion"`
	Close            bool          `json:"close"`
}
type ImagenUsuario struct {
	Imagen string `json:"uri"`
}

type Permisos struct {
	ModuloAutorizacion bool `json:"modulo_autorizacion"`
}

type HistoriasUsuarioNotificaciones struct {
	Usuarios       []*UsuarioHistoria `json:"usuarios"`
	Notificaciones *Notificaciones    `json:"notificaciones"`
	Permisos       *Permisos          `json:"permisos"`
	ModulosEtapa   *models.Etapa      `json:"etapa"`
}

func GetUsuariosHistoria(c *gin.Context) {
	idFiel := c.GetInt("id_residente")
	fechaActual := time.Now().In(tiempo.Local)
	hists := []*models.HistoriaFiel{}
	err := models.Db.Order("created_at desc").Where("created_at < ?", fechaActual).Where("fecha_fin > ?", fechaActual).Preload("Fiel", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "usuario_id")
	}).Preload("Fiel.Usuario", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id,imagen,nombre,usuario")
	}).Find(&hists).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener usuarios"), nil, c, http.StatusInternalServerError)
		return
	}
	residentes := []*UsuarioHistoria{}
	for _, hist := range hists {
		flag := true
		for _, res := range residentes {
			if res.ID == hist.FielID {
				flag = false
				break
			}

		}
		if flag {
			imagen := ""
			if hist.Fiel.Usuario.Imagen == "" {
				imagen = utils.DefaultUser
			} else {
				imagen = utils.SERVIMG + hist.Fiel.Usuario.Imagen
			}

			res := &UsuarioHistoria{ID: hist.Fiel.ID, Key: hist.Fiel.Usuario.Nombre, Nombre: hist.Fiel.Usuario.Usuario, Leido: false, Source: ImagenUsuario{Imagen: imagen}, FechaPublicacion: hist.CreatedAt}
			residentes = append(residentes, res)
		}
		var count int64
		err = models.Db.Model(&models.LecturaHistoria{}).Where("residente_id = ?", idFiel).Where("historia_residente_id = ?", hist.ID).Count(&count).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener usuarios"), nil, c, http.StatusInternalServerError)
			return
		}
		if count == 0 {
			for _, res := range residentes {
				if res.ID == hist.FielID {
					res.Leido = true
					break
				}

			}
		}

	}
	if len(residentes) == 1 {
		residentes = append(residentes, &UsuarioHistoria{Leido: false, ID: 10110101, Key: "Practical", Nombre: "Practical", Source: ImagenUsuario{Imagen: utils.LogoPractical}})
	}
	slice.Sort(residentes[:], func(i, j int) bool {
		return residentes[i].ID == uint(idFiel)
	})

	utils.CrearRespuesta(nil, residentes, c, http.StatusOK)

}

func GetUsuariosHistoriaNotifiaciones(c *gin.Context) {
	idFiel := c.GetInt("id_residente")
	idCasa := c.GetInt("id_casa")
	fechaActual := time.Now().In(tiempo.Local)
	hists := []*models.HistoriaFiel{}
	err := models.Db.Order("created_at desc").Where("created_at < ?", fechaActual).Where("fecha_fin > ?", fechaActual).Preload("Fiel", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "usuario_id")
	}).Preload("Fiel.Usuario", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id,imagen,nombre,usuario")
	}).Find(&hists).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener usuarios"), nil, c, http.StatusInternalServerError)
		return
	}
	residentes := []*UsuarioHistoria{}
	for _, hist := range hists {
		flag := true
		for _, res := range residentes {
			if res.ID == hist.FielID {
				flag = false
				break
			}

		}
		if flag {
			imagen := ""
			if hist.Fiel.Usuario.Imagen == "" {
				imagen = utils.DefaultUser
			} else {
				imagen = utils.SERVIMG + hist.Fiel.Usuario.Imagen
			}
			res := &UsuarioHistoria{ID: hist.Fiel.ID, Key: hist.Fiel.Usuario.Nombre, Nombre: hist.Fiel.Usuario.Usuario, Leido: false, Source: ImagenUsuario{Imagen: imagen}, FechaPublicacion: hist.CreatedAt, Close: false}
			residentes = append(residentes, res)
		}
		var count int64
		err = models.Db.Model(&models.LecturaHistoria{}).Where("residente_id = ?", idFiel).Where("historia_residente_id = ?", hist.ID).Count(&count).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener usuarios"), nil, c, http.StatusInternalServerError)
			return
		}
		if count == 0 {
			for _, res := range residentes {
				if res.ID == hist.FielID {
					res.Leido = true
					break
				}

			}
		}

	}

	slice.Sort(residentes[:], func(i int, j int) bool {
		return residentes[i].FechaPublicacion.After(residentes[j].FechaPublicacion)
	})
	slice.Sort(residentes[:], func(i, j int) bool {
		return residentes[i].Leido
	})

	slice.Sort(residentes[:], func(i, j int) bool {
		return residentes[i].ID == uint(idFiel)
	})

	if len(residentes) == 0 {
		residentes = append(residentes, &UsuarioHistoria{Leido: false, Key: "Practical", ID: 10110101, Nombre: "Practical", Source: ImagenUsuario{Imagen: utils.LogoPractical}})
	}
	for i := 0; i < len(residentes); i++ {
		if i+1 < len(residentes) {
			if !residentes[i+1].Leido {
				residentes[i].Close = true
			}

		} else {
			residentes[i].Close = true
		}

	}

	idParroquia := c.GetInt("id_etapa")
	notificaciones, err := obtenerNotificaciones(idFiel, idCasa, idParroquia)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al obtener usuarios"), nil, c, http.StatusInternalServerError)
		return
	}
	residente := &models.Fiel{}
	etapa := &models.Etapa{}
	err = models.Db.Select("autorizacion").First(&residente, idFiel).Error
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al obtener usuarios"), nil, c, http.StatusInternalServerError)
		return
	}
	err = models.Db.Select("pagos_tarjeta", "modulo_market", "modulo_publicacion", "modulo_votacion", "modulo_area_social", "modulo_equipo", "modulo_historia", "modulo_bitacora", "urbanizacion", "formulario_entrada", "formulario_salida", "modulo_alicuota", "modulo_emprendimiento", "modulo_camaras", "modulo_directiva", "modulo_galeria", "modulo_horarios", "modulo_mi_registro").First(&etapa, idParroquia).Error
	permisos := &Permisos{ModuloAutorizacion: residente.Autorizacion}
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al obtener usuarios"), nil, c, http.StatusInternalServerError)
		return
	}
	historias := &HistoriasUsuarioNotificaciones{Usuarios: residentes, Notificaciones: notificaciones, Permisos: permisos, ModulosEtapa: etapa}

	utils.CrearRespuesta(nil, historias, c, http.StatusOK)
}

type Historia struct {
	ID               uint      `json:"ID"`
	Contenido        string    `json:"content"`
	Tipo             string    `json:"type"`
	Leido            int       `json:"finish"`
	Views            *int64    `json:"views,omitempty"`
	IsUser           bool      `json:"is_user"`
	FechaPublicacion time.Time `json:"fecha_publicacion"`
}

func GetHistoriasDeUsuario(c *gin.Context) {
	idFielToken := uint(c.GetInt("id_residente"))
	idFiel := c.Param("id")

	historias := []*Historia{}
	if idFiel == "10110101" {
		var views int64
		views = 0
		historias = append(historias, &Historia{ID: 10110101, Contenido: utils.RutaTutorial, Tipo: "video", Leido: 1, Views: &views, IsUser: false})
		utils.CrearRespuesta(nil, historias, c, http.StatusOK)
		return
	}
	fechaActual := time.Now().In(tiempo.Local)
	hists := []*models.HistoriaFiel{}
	err := models.Db.Select("id", "is_video", "fecha_fin", "created_at", "url").Order("created_at asc").Where("created_at < ?", fechaActual).Where("fecha_fin > ?", fechaActual).Where("residente_id = ?", idFiel).Find(&hists).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener historias"), nil, c, http.StatusInternalServerError)
		return
	}

	for _, hist := range hists {
		historia := &Historia{}
		var count int64
		err = models.Db.Model(&models.LecturaHistoria{}).Where("historia_residente_id = ?", hist.ID).Where("residente_id = ?", idFielToken).Count(&count).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener historias"), nil, c, http.StatusInternalServerError)
			return
		}
		historia.ID = hist.ID
		historia.Contenido = utils.SERVIMG + hist.Url
		if hist.IsVideo {
			historia.Tipo = "video"
		} else {
			historia.Tipo = "image"
		}
		if count > 0 {
			historia.Leido = 1
		} else {
			historia.Leido = 0
		}
		idFielFormat, err := strconv.Atoi(idFiel)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error en parametros de peticion"), nil, c, http.StatusBadRequest)
			return
		}
		if uint(idFielFormat) == idFielToken {

			var visualizaciones int64
			err = models.Db.Model(&models.LecturaHistoria{}).Where("historia_residente_id = ?", hist.ID).Count(&visualizaciones).Error
			if err != nil {
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al obtener historias"), nil, c, http.StatusInternalServerError)
				return
			}
			historia.Views = &visualizaciones
			historia.IsUser = true
		}
		historia.FechaPublicacion = hist.CreatedAt
		historias = append(historias, historia)
	}
	// Resetear leidos si ya todas fueron vistas
	todasVistas := true
	for _, hist := range historias {
		if hist.Leido == 0 {
			todasVistas = false
		}
	}
	if todasVistas {
		for _, hist := range historias {
			hist.Leido = 0
		}
	}

	utils.CrearRespuesta(nil, historias, c, http.StatusOK)
}

func ConfirmarLecturaHistoria(c *gin.Context) {
	idHistoria := c.Param("id")
	idFiel := uint(c.GetInt("id_residente"))
	lectura := &models.LecturaHistoria{}
	idHist, err := strconv.Atoi(idHistoria)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al enviar parametros"), nil, c, http.StatusBadRequest)
		return

	}
	lectura.HistoriaFielID = uint(idHist)
	lectura.FielID = idFiel
	err = models.Db.Create(lectura).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al confirmar lectura"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Lectura confirmada", c, http.StatusOK)
}

func DeleteHistoria(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.HistoriaFiel{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar historia"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Historia eliminada exitosamente", c, http.StatusOK)
}
