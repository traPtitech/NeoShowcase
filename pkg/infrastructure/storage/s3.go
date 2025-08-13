package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

// S3Storage AmazonS3ストレージ
type S3Storage struct {
	client *s3.Client
	bucket string
}

// NewS3Storage S3Storageを生成する。指定したBucketはすでに存在している必要がある。
func NewS3Storage(bucket, accessKey, accessSecret, region, endpoint string) (*S3Storage, error) {
	creds := credentials.NewStaticCredentialsProvider(accessKey, accessSecret, "")
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load AWS config")
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	s := S3Storage{
		client: client,
		bucket: bucket,
	}
	return &s, nil
}

// Save ファイルをアップロードする
func (s3s *S3Storage) Save(filename string, src io.Reader) error {
	uploader := s3manager.NewUploader(s3s.client)
	_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s3s.bucket),
		Key:    aws.String(filename),
		Body:   src,
	})
	return err
}

// Open ファイルを開く
func (s3s *S3Storage) Open(filename string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s3s.bucket),
		Key:    aws.String(filename),
	}
	result, err := s3s.client.GetObject(context.Background(), input)
	if err != nil {
		return nil, domain.ErrFileNotFound
	}
	return result.Body, nil
}

// Delete ファイルを削除する
func (s3s *S3Storage) Delete(filename string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s3s.bucket),
		Key:    aws.String(filename),
	}
	_, err := s3s.client.DeleteObject(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}
