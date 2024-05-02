package usecases

import (
	"claim_api_demo/internal/pkg/domain/dtos"
	"claim_api_demo/internal/pkg/domain/entities"
	"claim_api_demo/internal/pkg/repositories"
	externalRepo "claim_api_demo/pkg/repositories"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

type ClaimUsecaseInterface interface {
	CreateClaim(cxt context.Context, req dtos.CreateClaimRequest, authUser entities.User) (dtos.Claim, error)
	GetList(ctx context.Context, authUser entities.User) ([]dtos.Claim, error)
	GetSpecific(ctx context.Context, authUser entities.User, claimID string) (dtos.Claim, error)
}

type ClaimUsecase struct {
	repo            repositories.ClaimRepositoryInterface
	fireStorageRepo repositories.FileStorageRepositoryInterface
	S3              externalRepo.S3RepositoryInterface
}

func (u *ClaimUsecase) CreateClaim(ctx context.Context, req dtos.CreateClaimRequest, authUser entities.User) (dtos.Claim, error) {
	bucket := os.Getenv("DOCUMENT_BUCKET")

	documents := []dtos.UploadFileResponse{}

	fileStorages, err := u.fireStorageRepo.GetByConditions(ctx, bson.D{{"id", bson.D{{"$in", req.StorageIDs}}}})
	if err != nil {
		return dtos.Claim{}, err
	}

	claim, err := u.repo.Create(ctx, entities.Claim{
		Ammount:   req.Ammount,
		UserID:    authUser.ID,
		Documents: fileStorages,
	})
	if err != nil {
		return dtos.Claim{}, err
	}

	for _, fileStorage := range fileStorages {
		url, err := u.S3.GeneratePresignedUrl(bucket, fileStorage.ID)
		if err != nil {
			return dtos.Claim{}, err
		}

		documents = append(documents, dtos.UploadFileResponse{
			ID:       fileStorage.ID,
			Filename: fileStorage.Filename,
			Url:      url,
		})
	}

	return dtos.Claim{
		ID:        claim.ID,
		Ammount:   claim.Ammount,
		Documents: documents,
	}, nil
}

func (u *ClaimUsecase) GetList(ctx context.Context, authUser entities.User) ([]dtos.Claim, error) {
	bucket := os.Getenv("DOCUMENT_BUCKET")
	claims, err := u.repo.GetByConditions(ctx, bson.D{{"user_id", authUser.ID}})
	if err != nil {
		return []dtos.Claim{}, err
	}

	var response []dtos.Claim
	for _, claim := range claims {
		item := dtos.Claim{
			ID:      claim.ID,
			Ammount: claim.Ammount,
			UserID:  claim.UserID,
		}
		for _, document := range claim.Documents {
			url, err := u.S3.GeneratePresignedUrl(bucket, document.ID)
			if err != nil {
				return []dtos.Claim{}, err
			}
			item.Documents = append(item.Documents, dtos.UploadFileResponse{
				ID:       document.ID,
				Filename: document.Filename,
				Url:      url,
			})

		}
		response = append(response, item)
	}

	return response, nil
}

func (u *ClaimUsecase) GetSpecific(ctx context.Context, authUser entities.User, claimID string) (dtos.Claim, error) {
	bucket := os.Getenv("DOCUMENT_BUCKET")
	claim, err := u.repo.TakeByConditions(ctx, bson.D{{"id", claimID}, {"user_id", authUser.ID}})
	if err != nil {
		return dtos.Claim{}, err
	}

	var response dtos.Claim
	for _, document := range claim.Documents {
		url, err := u.S3.GeneratePresignedUrl(bucket, document.ID)
		if err != nil {
			return dtos.Claim{}, err
		}
		response.Documents = append(response.Documents, dtos.UploadFileResponse{
			ID:       document.ID,
			Filename: document.Filename,
			Url:      url,
		})
	}
	return response, nil
}

func NewClaimUsecase(
	claimRepo repositories.ClaimRepositoryInterface,
	fireStorageRepo repositories.FileStorageRepositoryInterface,
	S3 externalRepo.S3RepositoryInterface,
) ClaimUsecaseInterface {
	return &ClaimUsecase{
		repo:            claimRepo,
		fireStorageRepo: fireStorageRepo,
		S3:              S3,
	}
}
