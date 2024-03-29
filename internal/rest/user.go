package rest

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
)

const (
	bucketName = "avatars"
)

func (h *Rest) getUsers(c *gin.Context) {
	result, err := h.services.GetUsers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Rest) getUserById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid is param")
		return
	}
	result, err := h.services.GetUserById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Rest) getRoles(c *gin.Context) {
	result, err := h.services.GetRoles()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Rest) createUser(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.services.CheckLogin(input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.User.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

func (h *Rest) deleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid is param")
		return
	}

	result, err := h.services.DeleteUser(c, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)

}

func (h *Rest) uploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	userId := c.Param("id")

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "No file is received")
		return
	}

	file, err := fileHeader.Open()
	defer file.Close()

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Cant open file")
		return
	}

	idFile := uuid.New().String()
	extension := filepath.Ext(fileHeader.Filename)

	newFileName := idFile + extension

	fileSize := fileHeader.Size

	info, mimeType, uploadErr := h.services.UploadFile(c, file, bucketName, newFileName, fileSize)

	if uploadErr != nil {
		newErrorResponse(c, http.StatusInternalServerError, uploadErr.Error())
	}

	input := models.Avatar{
		MimeType:   mimeType,
		BucketName: bucketName,
		FileName:   newFileName,
	}

	_, userImgErr := h.services.UpdateUserImg(c, userId, input)
	if userImgErr != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, info.Key)
}
