package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	requireRole         = "requireRole"
	duplicateCheck      = "duplicateCheck"
)

// Идентификация пользователя
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

// Проверка роли
func (h *Handler) setRoleUser(c *gin.Context) {
	c.Set(requireRole, "user")
}

func (h *Handler) setRoleExecutor(c *gin.Context) {
	c.Set(requireRole, "executor")
}

func (h *Handler) setRoleAdmin(c *gin.Context) {
	c.Set(requireRole, "admin")
}

func (h *Handler) checkRights(c *gin.Context) {
	role, ok := c.Get(requireRole)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "role not found")
		return
	}

	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	_, err := h.services.CheckRights(id, role)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}
}
