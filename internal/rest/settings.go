package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
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
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, result)
}
