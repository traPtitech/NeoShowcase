package storage

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Storage struct {
	Sess   *session.Session
	Bucket string
	Key    string
}

func (s3s *S3Storage) Save(filename string, src io.Reader) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(s3s.Sess)
	defer file.Close()
	_, err = io.Copy(file, src)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3s.Bucket),
		Key:    aws.String(s3s.Key),
		Body:   file,
	})
	return err
}
func (s3s *S3Storage) Open(filename string) (io.ReadCloser, error) {
	svc := s3.New(s3s.Sess)
	input := &s3.GetObjectInput{
		Bucket: aws.String(s3s.Bucket),
		Key:    aws.String(s3s.Key),
	}
	result, err := svc.GetObject(input)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (s3s *S3Storage) Delete(filename string) error {
	svc := s3.New(s3s.Sess)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s3s.Bucket),
		Key:    aws.String(s3s.Key),
	}
	_, err := svc.DeleteObject(input)
	if err != nil {
		return err
	}
	return nil
}

func (s3s *S3Storage) Move(sourcePath, destPath string) error {
	// Move LocalDir to Swift Storage
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %w", err)
	}
	uploader := s3manager.NewUploader(s3s.Sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3s.Bucket),
		Key:    aws.String(s3s.Key),
		Body:   inputFile,
	})
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %w", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %w", err)
	}
	return nil
}
