package repositories

import (
	"claim_api_demo/internal/pkg/domain/entities"
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileStorageRepositoryInterface interface {
	TakeByConditions(ctx context.Context, condition bson.D) (entities.FileStorage, error)
	Create(ctx context.Context, data entities.FileStorage) (entities.FileStorage, error)
	GetByConditions(ctx context.Context, conditions bson.D) ([]entities.FileStorage, error)
}

type FileStorageRepository struct {
	DB *mongo.Database
}

func (r *FileStorageRepository) TakeByConditions(ctx context.Context, conditions bson.D) (entities.FileStorage, error) {
	fileStorage := entities.FileStorage{}
	err := r.DB.Collection("file_storages").FindOne(ctx, conditions).Decode(&fileStorage)
	return fileStorage, err
}

func (r *FileStorageRepository) Create(ctx context.Context, fileStorage entities.FileStorage) (entities.FileStorage, error) {
	fileStorage.ID = uuid.New().String()
	_, err := r.DB.Collection("file_storages").InsertOne(ctx, fileStorage)
	return fileStorage, err
}

func (r *FileStorageRepository) GetByConditions(ctx context.Context, conditions bson.D) ([]entities.FileStorage, error) {
	var fileStorages []entities.FileStorage
	cursor, err := r.DB.Collection("file_storages").Find(ctx, conditions)
	if err != nil {
		return fileStorages, err
	}
	err = cursor.All(ctx, &fileStorages)
	return fileStorages, err
}

func NewFileStorageRepository(db *mongo.Database) FileStorageRepositoryInterface {
	return &FileStorageRepository{
		DB: db,
	}
}
