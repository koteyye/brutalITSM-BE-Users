package rest

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"github.com/koteyye/brutalITSM-BE-Users/internal/service"
	"net/http"
)

func (s *Rest) addSettings(c *gin.Context) {
	var input []models.Settings

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := s.services.AddSettings(input)
	if err != nil {
		if errors.Is(err, service.AllDuplicate) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (s *Rest) editSettings(c *gin.Context) {
	var input []models.Settings

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := s.services.EditSettings(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Rest) deleteSettings(c *gin.Context) {

	setObj := c.Query("settingsObject")
	if setObj == "" {
		newErrorResponse(c, http.StatusBadRequest, "не указан объект удаления")
		return
	}

	ids := c.QueryArray("ids")
	if len(ids) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "не указаны ID удаляемых записей")
		return
	}
	err := s.services.DeleteSettings(ids, setObj)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			newErrorResponse(c, http.StatusBadRequest, "по указанным ID не найдено записи для удаления")
			return
		case service.NoSettingsObject:
			newErrorResponse(c, http.StatusBadRequest, service.NoSettingsObject.Error())
			return
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"objectSettings": setObj,
		"ids":            ids,
		"result":         "Все удалено нах",
	})
}
