package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Rest) searchJob(c *gin.Context) {
	searchRequest := c.Param("jobName")
	if searchRequest == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid is param")
		return
	}
	searchResult, err := h.services.SearchJob(searchRequest)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, searchResult)
}

func (h *Rest) searchOrg(c *gin.Context) {
	searchRequest := c.Param("orgName")
	if searchRequest == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid is param")
		return
	}
	searchResult, err := h.services.SearchOrg(searchRequest)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, searchResult)
}
