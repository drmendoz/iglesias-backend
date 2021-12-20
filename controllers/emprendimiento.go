package controllers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/drmendoz/iglesias-backend/utils/tiempo"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm"
)

type EmprendimientoResponse struct {
	Recomendados []*models.Emprendimiento `json:"recomendados"`
	Cercas       []*models.Emprendimiento `json:"cercas"`
}

var p = message.NewPrinter(language.English)

func CreateEmprendimiento(c *gin.Context) {
	idFiel := uint(c.GetInt("id_residente"))
	idParroquia := uint(c.GetInt("id_etapa"))
	item := &models.Emprendimiento{}
	err := c.ShouldBindJSON(item)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error en formato de solicitud"), nil, c, http.StatusBadRequest)
		return
	}
	item.EmprendimientoImagenes = []*models.EmprendimientoImagen{}
	for _, imagen := range item.Imagenes {
		empImagen := &models.EmprendimientoImagen{}
		empImagen.Imagen = imagen
		item.EmprendimientoImagenes = append(item.EmprendimientoImagenes, empImagen)
	}
	// suscrito, err := verificarSuscripcion(idFiel)
	// if err != nil {
	// 	_ = c.Error(err)
	// 	utils.CrearRespuesta(errors.New("Error al crear emprendimiento"), nil, c, http.StatusInternalServerError)
	// 	return
	// }
	fechaActual := time.Now().In(tiempo.Local)
	fechaFin := fechaActual.AddDate(0, 1, 0)
	var numEmp int64
	err = models.Db.Model(&models.Emprendimiento{}).Where("fecha_publicacion < ?", fechaActual).Where("fecha_vencimiento > ?", fechaActual).Where("residente_id = ?", idFiel).Count(&numEmp).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear emprendimiento"), nil, c, http.StatusInternalServerError)
		return
	}
	//if suscrito {
	if numEmp == int64(utils.NumeroMaximoEmprendimiento) {
		utils.CrearRespuesta(errors.New("Ha alcanzado el limite de emprendimientos mensuales"), nil, c, http.StatusNotAcceptable)
		return

	}
	//} else {
	//	if numEmp == 1 {
	//		utils.CrearRespuesta(errors.New("Ha alcanzado el limite de emprendimientos mensuales"), nil, c, http.StatusNotAcceptable)
	//		return

	//	}
	//}
	if item.PrecioLabel != "" {
		item.Precio, err = strconv.ParseFloat(item.PrecioLabel, 64)
		if err != nil {
			utils.CrearRespuesta(errors.New("Error en formato de precio"), nil, c, http.StatusBadRequest)
			return
		}
	}
	item.Premium = true
	item.FechaPublicacion = fechaActual
	item.FechaVencimiento = fechaFin
	item.ParroquiaID = idParroquia
	item.FielID = idFiel
	//count := 0
	// for _, imagen := range item.Imagenes {
	// 	empImagen := &models.EmprendimientoImagen{}
	// 	if imagen != "" {
	// 		ct := fmt.Sprintf("%d", count)
	// 		imagen, err = img.FromBase64ToImage(imagen, "emprendimientos/"+time.Now().Format(time.RFC3339)+ct)
	// 		if err != nil {
	// 			_ = c.Error(err)
	// 			utils.CrearRespuesta(errors.New("Error en formato de imagen"), nil, c, http.StatusOK)
	// 			return
	// 		}
	// 	}
	// 	count++
	// 	empImagen.Imagen = imagen
	// 	item.EmprendimientoImagenes = append(item.EmprendimientoImagenes, empImagen)
	// }
	err = models.Db.Create(item).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear item"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "Emprendimiento creado correctamente", c, http.StatusCreated)

}

func UpdateEmprendimiento(c *gin.Context) {
	item := &models.Emprendimiento{}

	err := c.ShouldBindJSON(item)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Where("id = ?", id).Updates(item).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar item"), nil, c, http.StatusInternalServerError)
		return
	}
	if item.Imagen != "" {
		idUrb := fmt.Sprintf("%d", item.ID)
		item.Imagen, err = img.FromBase64ToImage(item.Imagen, "items/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear item "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Emprendimiento{}).Where("id = ?", item.ID).Update("imagen", item.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar item"), nil, c, http.StatusInternalServerError)
			return
		}
		item.Imagen = utils.SERVIMG + item.Imagen

	} else {
		item.Imagen = utils.SERVIMG + "default_user.png"
	}
	utils.CrearRespuesta(err, "Item actualizado correctamente", c, http.StatusOK)
}

func DeleteEmprendimiento(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Emprendimiento{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar emprendimiento"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Emprendimiento eliminado exitosamente", c, http.StatusOK)
}

func ObtenerEmprendimientos(c *gin.Context) {
	idParroquia := c.GetInt("id_etapa")
	filtro := c.Query("filtro")
	categoria := c.Query("id_categoria")
	idCat := 0
	var err error
	if categoria != "" {
		idCat, err = strconv.Atoi(categoria)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error en id de categoria"), nil, c, http.StatusBadRequest)
			return
		}
	}

	emps := &EmprendimientoResponse{}
	emps.Cercas = []*models.Emprendimiento{}
	emps.Recomendados = []*models.Emprendimiento{}
	fechaActual := time.Now().In(tiempo.Local)

	//Recomendados
	err = models.Db.Where("estado = ?", "VIG").Where(&models.Emprendimiento{CategoriaMarketID: uint(idCat), Premium: true}).Where("fecha_publicacion < ?", fechaActual).Where("titulo like ?", "%"+filtro+"%").Preload("EmprendimientoImagenes").Preload("Fiel.Usuario", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "imagen", "telefono", "usuario", "telefono")
	}).Order("created_at desc").Find(&emps.Recomendados).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener emprendimientos"), nil, c, http.StatusInternalServerError)
		return
	}

	//Cercas
	err = models.Db.Where("estado = ?", "VIG").Order("premium ASC").Where(&models.Emprendimiento{CategoriaMarketID: uint(idCat), ParroquiaID: uint(idParroquia)}).Where("fecha_publicacion < ?", fechaActual).Where("titulo like ?", "%"+filtro+"%").Preload("EmprendimientoImagenes").Preload("Fiel.Usuario", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "imagen", "telefono", "usuario")
	}).Order("created_at desc").Find(&emps.Cercas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener emprendimientos"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, emp := range emps.Cercas {
		// re := regexp.MustCompile(`\.`)
		// emp.Descripcion = re.ReplaceAllString(emp.Descripcion, ".\n")
		emp.PrecioLabel = p.Sprintf("%.0f", emp.Precio)
		if emp.EmprendimientoImagenes == nil {
			emp.Imagen = utils.DefaultAreaSocial
		} else {
			if len(emp.EmprendimientoImagenes) > 0 {
				if emp.EmprendimientoImagenes[0].Imagen == "" {
					emp.Imagen = utils.DefaultAreaSocial
				} else {

					emp.Imagen = emp.EmprendimientoImagenes[0].Imagen
				}
			} else {
				emp.Imagen = utils.DefaultAreaSocial
			}

		}
		if emp.Fiel.Usuario.Imagen != "" {
			if !strings.HasPrefix(emp.Fiel.Usuario.Imagen, "https://") {
				emp.ImagenUsuario = utils.SERVIMG + emp.Fiel.Usuario.Imagen
			}
		}
		emp.TelefonoUsuario = emp.Fiel.Usuario.Telefono
		emp.NombreUsuario = emp.Fiel.Usuario.Usuario
		for _, img := range emp.EmprendimientoImagenes {
			if img.Imagen == "" {
				img.Imagen = utils.DefaultAreaSocial
			}
			emp.Imagenes = append(emp.Imagenes, img.Imagen)
		}

	}
	for _, emp := range emps.Recomendados {
		// re := regexp.MustCompile(`\.`)
		// emp.Descripcion = re.ReplaceAllString(emp.Descripcion, ".\n")
		emp.PrecioLabel = p.Sprintf("%.0f", emp.Precio)
		if emp.EmprendimientoImagenes == nil {
			emp.Imagen = utils.DefaultAreaSocial
		} else {
			if len(emp.EmprendimientoImagenes) > 0 {
				if emp.EmprendimientoImagenes[0].Imagen == "" {
					emp.Imagen = utils.DefaultAreaSocial
				} else {

					emp.Imagen = emp.EmprendimientoImagenes[0].Imagen
				}
			} else {

				emp.Imagen = utils.DefaultAreaSocial
			}
		}
		if emp.Fiel.Usuario.Imagen != "" {
			if !strings.HasPrefix(emp.Fiel.Usuario.Imagen, "https://") {
				emp.ImagenUsuario = utils.SERVIMG + emp.Fiel.Usuario.Imagen
			}
		}
		emp.TelefonoUsuario = emp.Fiel.Usuario.Telefono
		emp.NombreUsuario = emp.Fiel.Usuario.Usuario
		for _, img := range emp.EmprendimientoImagenes {
			if img.Imagen == "" {
				img.Imagen = utils.DefaultAreaSocial
			}
			emp.Imagenes = append(emp.Imagenes, img.Imagen)
		}
		emp.EmprendimientoImagenes = nil
	}
	utils.CrearRespuesta(nil, emps, c, http.StatusOK)
}

func ObtenerEmprendimientosPorId(c *gin.Context) {
	idEmprendimiento := c.Param("id")
	emp := &models.Emprendimiento{}
	err := models.Db.Where("estado = 'VIG'").Preload("EmprendimientoImagenes").Preload("Fiel", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "usuario_id")
	}).Preload("Fiel.Usuario", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "usuario", "telefono", "imagen")
	}).First(&emp, idEmprendimiento).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Emprendimiento no encontrado"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener emprendimiento"), nil, c, http.StatusInternalServerError)
		return
	}
	emp.Imagenes = []string{}
	emp.PrecioLabel = p.Sprintf("%.0f", emp.Precio)
	for _, img := range emp.EmprendimientoImagenes {
		emp.Imagenes = append(emp.Imagenes, img.Imagen)
	}

	emp.Precio = math.Round(emp.Precio*100) / 100
	emp.TelefonoUsuario = emp.Fiel.Usuario.Celular
	emp.NombreUsuario = emp.Fiel.Usuario.Usuario
	if emp.Fiel.Usuario.Imagen != "" {
		if !strings.HasPrefix(emp.Fiel.Usuario.Imagen, "https://") {
			emp.ImagenUsuario = utils.SERVIMG + emp.Fiel.Usuario.Imagen
		}
	}
	emp.Fiel = nil
	utils.CrearRespuesta(nil, emp, c, http.StatusOK)
}

func ObtenerEmprendimientosUsuarios(c *gin.Context) {
	filtro := c.Query("filtro")
	idFiel := uint(c.GetInt("id_residente"))
	categoria := c.Query("id_categoria")
	idCat := 0
	var err error
	if categoria != "" {
		idCat, err = strconv.Atoi(categoria)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error en id de categoria"), nil, c, http.StatusBadRequest)
			return
		}
	}

	fechaActual := time.Now().In(tiempo.Local)
	emps := []*models.Emprendimiento{}

	err = models.Db.Where("estado = ?", "VIG").Where(&models.Emprendimiento{CategoriaMarketID: uint(idCat), FielID: idFiel}).Where("fecha_vencimiento > ?", fechaActual).Where("titulo like ?", "%"+filtro+"%").Preload("EmprendimientoImagenes").Order("created_at desc").Find(&emps).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener emprendimientos"), nil, c, http.StatusInternalServerError)
		return
	}

	for _, emp := range emps {

		emp.PrecioLabel = fmt.Sprintf("%.2f", emp.Precio)
		if emp.EmprendimientoImagenes != nil {
			emp.Imagen = utils.DefaultAreaSocial
		} else {
			if emp.EmprendimientoImagenes[0].Imagen == "" {
				emp.Imagen = utils.DefaultAreaSocial
			} else {

				emp.Imagen = emp.EmprendimientoImagenes[0].Imagen
			}
		}
		emp.Imagenes = []string{}
		for _, img := range emp.EmprendimientoImagenes {
			emp.Imagenes = append(emp.Imagenes, img.Imagen)
		}
		emp.EmprendimientoImagenes = nil
	}

	utils.CrearRespuesta(nil, emps, c, http.StatusOK)
}

type FielEmprendimiento struct {
	Fiel           *FielReduce              `json:"residente"`
	Emprendimiento []*models.Emprendimiento `json:"emprendimientos"`
}
type FielReduce struct {
	Nombre string `json:"nombre"`
	Imagen string `json:"imagen"`
}

func ObtenerEmprendimientoFiel(c *gin.Context) {
	idRes := c.Param("id")
	idFiel := 0
	var err error
	idFiel, err = strconv.Atoi(idRes)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error en formato de id"), nil, c, http.StatusBadRequest)
		return
	}
	res := &models.Fiel{}
	err = models.Db.Joins("Usuario").First(res, idFiel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Fiel no encontrado"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error encontrando residente"), nil, c, http.StatusInternalServerError)
		return
	}
	if res.Usuario.Imagen != "" {
		res.Usuario.Imagen = utils.SERVIMG + res.Usuario.Imagen
	} else {
		res.Usuario.Imagen = utils.DefaultUser
	}
	fechaActual := time.Now().In(tiempo.Local)
	emps := []*models.Emprendimiento{}
	err = models.Db.Where("estado = 'VIG'").Where("residente_id = ?", idFiel).Where("fecha_publicacion < ?", fechaActual).Where("fecha_vencimiento > ?", fechaActual).Preload("EmprendimientoImagenes").Order("created_at desc").Find(&emps).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Emprendimiento no encontrado"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener emprendimiento"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, emp := range emps {
		emp.Imagenes = []string{}

		emp.PrecioLabel = fmt.Sprintf("%.2f", emp.Precio)
		for _, img := range emp.EmprendimientoImagenes {
			emp.Imagenes = append(emp.Imagenes, img.Imagen)
		}
		emp.TelefonoUsuario = res.Usuario.Celular
		emp.NombreUsuario = res.Usuario.Nombre

		emp.Precio = math.Round(emp.Precio*100) / 100
		emp.Fiel = nil
	}

	resEmp := &FielEmprendimiento{Fiel: &FielReduce{Nombre: res.Usuario.Nombre, Imagen: res.Usuario.Imagen}, Emprendimiento: emps}

	utils.CrearRespuesta(nil, resEmp, c, http.StatusOK)
}

func ActualizarEmprendimiento(c *gin.Context) {
	idEmprendimiento := c.Param("id")
	id, err := strconv.Atoi(idEmprendimiento)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error de id "), nil, c, http.StatusBadRequest)
		return
	}
	item := &models.Emprendimiento{}
	err = c.ShouldBindJSON(item)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error en formato de solicitud"), nil, c, http.StatusBadRequest)
		return
	}

	tx := models.Db.Begin()
	err = tx.Where(" emprendimiento_id = ?", id).Delete(&models.EmprendimientoImagen{}).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar emprendimiento"), nil, c, http.StatusInternalServerError)
		return
	}
	item.EmprendimientoImagenes = []*models.EmprendimientoImagen{}
	for _, imagen := range item.Imagenes {
		empImagen := &models.EmprendimientoImagen{}
		empImagen.Imagen = imagen
		item.EmprendimientoImagenes = append(item.EmprendimientoImagenes, empImagen)
	}

	item.ID = uint(id)
	if item.PrecioLabel != "" {
		item.Precio, err = strconv.ParseFloat(item.PrecioLabel, 64)
		if err != nil {
			utils.CrearRespuesta(errors.New("Error en formato de precio"), nil, c, http.StatusBadRequest)
			return
		}
	}
	err = tx.Updates(item).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar emprendimiento"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Emprendimiento actualizado correctamente", c, http.StatusCreated)
}
