package storage

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Storage AmazonS3ストレージ
type S3Storage struct {
	sess   *session.Session
	bucket string
}

// NewS3Storage S3Storageを生成する。指定したBucketはすでに存在している必要がある。
func NewS3Storage(bucket, accessKey, accessSecret, region, endpoint string) (*S3Storage, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, accessSecret, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := session.NewSession(s3Config)
	if err != nil {
		return nil, fmt.Errorf("failed to new session: %w", err)
	}
	s := S3Storage{
		sess:   sess,
		bucket: bucket,
	}
	return &s, nil
}

// Save ファイルをアップロードする
func (s3s *S3Storage) Save(filename string, src io.Reader) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(s3s.sess)
	_, err = io.Copy(file, src)
	if err != nil {
		return err
	}
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3s.bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return err
	}
	return os.Remove(filename)
}

// Open ファイルを開く
func (s3s *S3Storage) Open(filename string) (io.ReadCloser, error) {
	svc := s3.New(s3s.sess)
	input := &s3.GetObjectInput{
		Bucket: aws.String(s3s.bucket),
		Key:    aws.String(filename),
	}
	result, err := svc.GetObject(input)
	if err != nil {
		return nil, ErrFileNotFound
	}
	return result.Body, nil
}

// Delete ファイルを削除する
func (s3s *S3Storage) Delete(filename string) error {
	svc := s3.New(s3s.sess)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s3s.bucket),
		Key:    aws.String(filename),
	}
	_, err := svc.DeleteObject(input)
	if err != nil {
		return err
	}
	return nil
}

// Move 指定したローカルのファイルをストレージへ移動する。destPathは使用されない。
func (s3s *S3Storage) Move(filename, destPath string) error {
	// Move LocalDir to Swift Storage
	inputFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %w", err)
	}
	uploader := s3manager.NewUploader(s3s.sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3s.bucket),
		Key:    aws.String(destPath),
		Body:   inputFile,
	})
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %w", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("failed removing original file: %w", err)
	}
	return nil
}
