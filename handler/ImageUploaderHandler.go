package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tipbk/sneakfeed-service/dto"
	"github.com/tipbk/sneakfeed-service/service"
	"github.com/tipbk/sneakfeed-service/util"
)

type imageUploaderHandler struct {
	imageUploaderService service.ImageUploaderService
}

type ImageUploaderHandler interface {
	UploadImage(c *gin.Context)
}

func NewImageUploaderHandler(imageUploaderService service.ImageUploaderService) ImageUploaderHandler {
	return &imageUploaderHandler{
		imageUploaderService: imageUploaderService,
	}
}

func (h *imageUploaderHandler) UploadImage(c *gin.Context) {
	var uploadRequest dto.UploadRequest
	if err := c.ShouldBindJSON(&uploadRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if uploadRequest.File == "" {
		c.JSON(http.StatusBadRequest, util.GenerateFailedResponse("cannot add empty file"))
		return
	}

	uploaderResponse, err := h.imageUploaderService.UploadImage(uploadRequest.File)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GenerateFailedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.GenerateSuccessResponse(uploaderResponse.Data.Url))
}
