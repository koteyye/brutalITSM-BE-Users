package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) getUsers(c *gin.Context) {
	//id, ok := c.Get(userCtx)
	//if !ok {
	//	newErrorResponse(c, http.StatusInternalServerError, "user id not found")
	//	return
	//}
	//
	//var input models.UserList
	//if err := c.BindJSON(&input); err != nil {
	//	newErrorResponse(c, http.StatusBadRequest, err.Error())
	//	return
	//}

}

func (h *Handler) getUserById(c *gin.Context) {

}
