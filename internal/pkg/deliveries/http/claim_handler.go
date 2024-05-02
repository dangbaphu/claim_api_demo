package http

import (
	"claim_api_demo/internal/pkg/domain/dtos"
	"claim_api_demo/internal/pkg/usecases"
	"claim_api_demo/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ClaimHandler struct {
	usecase usecases.ClaimUsecaseInterface
}

func (h *ClaimHandler) Create(c *gin.Context) {
	req := dtos.CreateClaimRequest{}
	err := c.ShouldBindWith(&req, binding.JSON)
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

	user := middleware.GetUserFromContext(c)
	claim, err := h.usecase.CreateClaim(c, req, user)
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
		Data: dtos.BaseResponse{
			Data: claim,
		},
	}
	c.JSON(http.StatusOK, data)
}

func (h *ClaimHandler) GetList(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	claims, err := h.usecase.GetList(c, user)
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
		Data: dtos.BaseResponse{
			Data: claims,
		},
	}
	c.JSON(http.StatusOK, data)
}

func (h *ClaimHandler) GetSpecific(c *gin.Context) {
	claimID := c.Param("claim_id")
	user := middleware.GetUserFromContext(c)
	claim, err := h.usecase.GetSpecific(c, user, claimID)
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
		Data: dtos.BaseResponse{
			Data: claim,
		},
	}
	c.JSON(http.StatusOK, data)
}

func NewClaimHandler(usecase usecases.ClaimUsecaseInterface) *ClaimHandler {
	return &ClaimHandler{
		usecase: usecase,
	}
}
