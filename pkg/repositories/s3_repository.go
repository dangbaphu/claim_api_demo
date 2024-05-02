package repositories

import (
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
)

type S3RepositoryInterface interface {
	UploadFile(bucket string, key string, file multipart.File) error
	GeneratePresignedUrl(bucket string, key string) (string, error)
}

type S3Repository struct {
	service    *s3.S3
	s3Uploader *s3manager.Uploader
}

func NewS3Repository(logger *logrus.Logger) (S3RepositoryInterface, error) {
	config := &aws.Config{
		Region: aws.String(os.Getenv("REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ID"),
			os.Getenv("AWS_SECRET"),
			"",
		),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	session, err := session.NewSession(config)
	if err != nil {
		logger.Error("cannot create s3 client")
	}
	return &S3Repository{
		service:    s3.New(session),
		s3Uploader: s3manager.NewUploader(session),
	}, err
}

func (r *S3Repository) UploadFile(bucket string, key string, file multipart.File) error {
	_, err := r.service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	return err
}

func (r *S3Repository) GeneratePresignedUrl(bucket string, key string) (string, error) {
	req, _ := r.service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(15 * time.Minute)
}
