package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"net/http"
)

func (h *Rest) addSettings(c *gin.Context) {
	var input models.Settings

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

}
