package database

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// DBMS postgres type
	DBMS = "postgres"
)

// Database struct
type Database struct {
	DB *mongo.Database
}

func NewDatabase(c context.Context, logger *logrus.Logger) (Database, error) {
	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		logger.Error("cannot create mongo client")
	}
	db := client.Database(dbName)
	return Database{db}, err
}
