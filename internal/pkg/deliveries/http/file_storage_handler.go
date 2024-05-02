package http

import (
	"claim_api_demo/internal/pkg/domain/dtos"
	"claim_api_demo/internal/pkg/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileStorageHandler struct {
	usecase usecases.FileStorageUsecaseInterface
}

// Login func
func (h *FileStorageHandler) UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		data := dtos.BaseResponse{
			Status: http.StatusBadRequest,
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}

	filename := fileHeader.Filename
	file, err := fileHeader.Open()
	if err != nil {
		data := dtos.BaseResponse{
			Status: http.StatusBadRequest,
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}

	fileStorage, err := h.usecase.UploadFile(c, file, filename)
	if err != nil {
		data := dtos.BaseResponse{
			Status: http.StatusNotFound,
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		}
		c.JSON(http.StatusBadRequest, data)
		return
	}

	data := dtos.BaseResponse{
		Status: http.StatusOK,
		Data:   fileStorage,
	}
	c.JSON(http.StatusOK, data)
}

func NewFileStorageHandler(usecase usecases.FileStorageUsecaseInterface) *FileStorageHandler {
	return &FileStorageHandler{
		usecase: usecase,
	}
}
