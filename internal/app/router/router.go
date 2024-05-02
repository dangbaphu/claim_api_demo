package router

import (
	handlers "claim_api_demo/internal/pkg/deliveries/http"
	"claim_api_demo/internal/pkg/domain/dtos"
	"claim_api_demo/internal/pkg/repositories"
	"claim_api_demo/internal/pkg/usecases"
	"claim_api_demo/pkg/middleware"
	externalRepo "claim_api_demo/pkg/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Router struct {
	Engine *gin.Engine
	DB     *mongo.Database
	S3     externalRepo.S3RepositoryInterface
}

func (r *Router) InitRoute() {
	r.Engine.Use(gin.Logger())
	r.Engine.Use(gin.Recovery())

	fr := repositories.NewFileStorageRepository(r.DB)

	// auth
	ur := repositories.NewUserRepository(r.DB)
	au := usecases.NewAuthUsecase(ur)
	ah := handlers.NewAuthHandler(au)

	// claim
	cr := repositories.NewClaimRepository(r.DB)
	cu := usecases.NewClaimUsecase(cr, fr, r.S3)
	ch := handlers.NewClaimHandler(cu)

	// file_storege
	fu := usecases.NewFileStorageUsecase(fr, r.S3)
	fh := handlers.NewFileStorageHandler(fu)
	r.Engine.GET("/health-check", func(c *gin.Context) {
		data := dtos.BaseResponse{
			Status: http.StatusOK,
			Data:   gin.H{"message": "Health check OK!"},
			Error:  nil,
		}
		c.JSON(http.StatusOK, data)
	})

	api := r.Engine.Group("/api")
	{
		api.POST("/login", ah.Login)
		api.Use(middleware.AuthUser(ur))
		// router api for todo
		claimAPI := api.Group("/claims")
		{
			claimAPI.POST("", ch.Create)
			claimAPI.GET("", ch.GetList)
			claimAPI.GET("/:claim_id", ch.GetSpecific)
		}
		fileStorage := api.Group("/file_storage")
		{
			fileStorage.POST("", fh.UploadFile)
		}
	}
}
