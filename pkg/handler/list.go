package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) getUsers(c *gin.Context) {
	id, _ := c.Get(userCtx)
	logrus.Info(userCtx)
	logrus.Info(id)
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getUserById(c *gin.Context) {

}
