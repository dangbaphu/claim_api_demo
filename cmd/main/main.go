package main

import (
	"claim_api_demo/internal/app/router"
	"claim_api_demo/pkg/database"
	"claim_api_demo/pkg/repositories"
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Post struct {
	Tile    string
	Content string
}

func main() {
	port := os.Getenv("PORT")
	log := logrus.New()
	c := context.Background()

	db, _ := database.NewDatabase(c, log)

	s3, _ := repositories.NewS3Repository(log)
	engine := gin.New()

	r := &router.Router{
		Engine: engine,
		DB:     db.DB,
		S3:     s3,
	}

	r.InitRoute()
	engine.Run(":" + port)

}
