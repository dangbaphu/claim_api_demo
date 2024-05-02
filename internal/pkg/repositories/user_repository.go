package repositories

import (
	"claim_api_demo/internal/pkg/domain/entities"
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface {
	TakeByConditions(ctx context.Context, condition bson.D) (entities.User, error)
	Create(ctx context.Context, data entities.User) (entities.User, error)
}

type UserRepository struct {
	DB *mongo.Database
}

func (r *UserRepository) TakeByConditions(ctx context.Context, conditions bson.D) (entities.User, error) {
	user := entities.User{}
	err := r.DB.Collection("users").FindOne(ctx, conditions).Decode(&user)
	return user, err
}

func (r *UserRepository) Create(ctx context.Context, user entities.User) (entities.User, error) {
	user.ID = uuid.New().String()
	user.Email = "baphu95@gmail.com"
	user.Password = "XohImNooBHFR0OVvjcYpJ3NgPQ1qq73WKhHvch0VQtg="
	_, err := r.DB.Collection("users").InsertOne(ctx, user)
	return user, err
}

func NewUserRepository(db *mongo.Database) UserRepositoryInterface {
	return &UserRepository{
		DB: db,
	}
}
