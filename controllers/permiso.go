package controllers

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPermisos(c *gin.Context) {
	idUsuario := c.GetInt("id_usuario")
	permisos := []*models.Permiso{}
	err := models.Db.Where("usuario_id = ?", idUsuario).Find(&permisos).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener permisos"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(err, permisos, c, http.StatusOK)
}

func GetPermisoPorId(c *gin.Context) {
	permiso := &models.Permiso{}
	id := c.Param("id")
	err := models.Db.First(permiso, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Permiso no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener permiso"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(nil, permiso, c, http.StatusOK)
}

func CreatePermiso(c *gin.Context) {
	permiso := &models.Permiso{}
	err := c.ShouldBindJSON(permiso)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	err = tx.Create(permiso).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear permiso"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Permiso creado correctamente", c, http.StatusCreated)

}

func UpdatePermiso(c *gin.Context) {
	permiso := &models.Permiso{}

	err := c.ShouldBindJSON(permiso)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Where("id = ?", id).Updates(permiso).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar permiso"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Permiso actualizada correctamente", c, http.StatusOK)
}

func DeletePermiso(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Permiso{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar permiso"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Permiso eliminada exitosamente", c, http.StatusOK)
}
