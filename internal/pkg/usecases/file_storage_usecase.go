package usecases

import (
	"claim_api_demo/internal/pkg/domain/dtos"
	"claim_api_demo/internal/pkg/domain/entities"
	"claim_api_demo/internal/pkg/repositories"
	externalRepo "claim_api_demo/pkg/repositories"
	"context"
	"mime/multipart"
	"os"
)

type FileStorageUsecaseInterface interface {
	UploadFile(cxt context.Context, file multipart.File, filename string) (dtos.UploadFileResponse, error)
}

type FileStorageUsecase struct {
	repo repositories.FileStorageRepositoryInterface
	S3   externalRepo.S3RepositoryInterface
}

func (u *FileStorageUsecase) UploadFile(ctx context.Context, file multipart.File, filename string) (dtos.UploadFileResponse, error) {
	bucket := os.Getenv("DOCUMENT_BUCKET")
	fileStorage, err := u.repo.Create(ctx, entities.FileStorage{
		Filename: filename,
	})
	if err != nil {
		return dtos.UploadFileResponse{}, err
	}
	err = u.S3.UploadFile(bucket, fileStorage.ID, file)
	if err != nil {
		return dtos.UploadFileResponse{}, err
	}
	url, err := u.S3.GeneratePresignedUrl(bucket, fileStorage.ID)
	if err != nil {
		return dtos.UploadFileResponse{}, err
	}

	return dtos.UploadFileResponse{
		ID:       fileStorage.ID,
		Filename: fileStorage.Filename,
		Url:      url,
	}, nil
}

func NewFileStorageUsecase(fileStorageRepo repositories.FileStorageRepositoryInterface, s3 externalRepo.S3RepositoryInterface) FileStorageUsecaseInterface {
	return &FileStorageUsecase{
		repo: fileStorageRepo,
		S3:   s3,
	}
}
