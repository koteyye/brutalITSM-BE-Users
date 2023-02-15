package handler

import (
	"brutalITSM-BE-Users/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {

	var message string
	var input models.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	message = fmt.Sprintf("Поздравляю! Создан глубоко уважаемый пользователь %s %s %s", input.Lastname, input.Firstname, input.Middlename)

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"message": message,
	})

}

func (h *Handler) signIn(c *gin.Context) {

}
