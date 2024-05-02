package repositories

import (
	"claim_api_demo/internal/pkg/domain/entities"
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository interface
type ClaimRepositoryInterface interface {
	TakeByConditions(ctx context.Context, condition bson.D) (entities.Claim, error)
	Create(ctx context.Context, data entities.Claim) (entities.Claim, error)
	GetByConditions(ctx context.Context, conditions bson.D) ([]entities.Claim, error)
}

type ClaimRepository struct {
	DB *mongo.Database
}

func (r *ClaimRepository) TakeByConditions(ctx context.Context, conditions bson.D) (entities.Claim, error) {
	claim := entities.Claim{}
	err := r.DB.Collection("claims").FindOne(ctx, conditions).Decode(&claim)
	return claim, err
}

func (r *ClaimRepository) Create(ctx context.Context, data entities.Claim) (entities.Claim, error) {
	data.ID = uuid.New().String()
	_, err := r.DB.Collection("claims").InsertOne(ctx, data)
	return data, err
}

func (r *ClaimRepository) GetByConditions(ctx context.Context, conditions bson.D) ([]entities.Claim, error) {
	var claims []entities.Claim
	cursor, err := r.DB.Collection("claims").Find(ctx, conditions)
	if err != nil {
		return claims, err
	}
	err = cursor.All(ctx, &claims)
	return claims, err
}

func NewClaimRepository(db *mongo.Database) ClaimRepositoryInterface {
	return &ClaimRepository{
		DB: db,
	}
}
