package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type signInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SignIn
// @Summary Sign in account
// @Description  Sign in your account with login and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input body signInInput true "account info"
// @Success      200  {object}  model.Account
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /accounts/{id} [get]

func (h *Rest) signIn(c *gin.Context) {

	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Login, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": token,
	})

}

func (h *Rest) me(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "user unauthorized")
		return
	}

	user, err := h.services.Authorization.Me(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, user)
}
