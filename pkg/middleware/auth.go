package middleware

import (
	"claim_api_demo/internal/pkg/domain/dtos"
	"claim_api_demo/internal/pkg/domain/entities"
	"claim_api_demo/internal/pkg/repositories"
	"claim_api_demo/pkg/auth"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

func AuthUser(userRepo repositories.UserRepositoryInterface) gin.HandlerFunc {
	secret := os.Getenv("ACCESS_SECRET")
	return func(c *gin.Context) {
		header := authHeader{}
		err := c.ShouldBindHeader(&header)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dtos.BaseResponse{
				Status: http.StatusUnauthorized,
				Error: &dtos.ErrorResponse{
					ErrorMessage: err.Error(),
				},
			})
			c.Abort()
			return
		}

		idTokenHeader := strings.Split(header.IDToken, "Bearer ")

		if len(idTokenHeader) < 2 {
			err := errors.New("must provide Authorization header with format `Bearer {token}`")
			c.JSON(http.StatusUnauthorized, dtos.BaseResponse{
				Status: http.StatusUnauthorized,
				Error: &dtos.ErrorResponse{
					ErrorMessage: err.Error(),
				},
			})
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(idTokenHeader[1], secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dtos.BaseResponse{
				Status: http.StatusUnauthorized,
				Error: &dtos.ErrorResponse{
					ErrorMessage: err.Error(),
				},
			})
			c.Abort()
			return
		}

		userID := claims["user_id"]
		user, err := userRepo.TakeByConditions(c, bson.D{{"id", userID}})
		if err != nil {
			c.JSON(http.StatusUnauthorized, dtos.BaseResponse{
				Status: http.StatusUnauthorized,
				Error: &dtos.ErrorResponse{
					ErrorMessage: err.Error(),
				},
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// GetUserFromContext func
func GetUserFromContext(c *gin.Context) entities.User {
	value, exist := c.Get("user")
	if !exist {
		return entities.User{}
	}
	return value.(entities.User)
}
