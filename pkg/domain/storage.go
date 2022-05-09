package domain

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrFileNotFound ファイルが存在しない
	ErrFileNotFound = errors.New("not found")
)

// Storage ストレージインターフェース
type Storage interface {
	Save(filename string, src io.Reader) error
	Open(filename string) (io.ReadCloser, error)
	Delete(filename string) error
	Move(filename, destPath string) error // LocalFile to Storage
}

// SaveArtifact Artifactをtar形式で保存する
func SaveArtifact(s Storage, filename string, dstpath string) error {
	if err := s.Move(filename, dstpath); err != nil {
		return fmt.Errorf("failed to save artifact tar file: %w", err)
	}

	return nil
}

// SaveLogFile ログファイルを保存する
func SaveLogFile(s Storage, filename string, dstpath string, buildid string) error {
	if err := s.Move(filename, dstpath); err != nil {
		return fmt.Errorf("failed to move build log: %w", err)
	}
	return nil
}

// ExtractTarToDir tarファイルをディレクトリに展開する
func ExtractTarToDir(s Storage, sourcePath, destPath string) error {
	inputFile, err := s.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %w", err)
	}
	defer inputFile.Close()

	tr := tar.NewReader(inputFile)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("bad tar file: %w", err)
		}

		path := filepath.Join(destPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, header.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), header.FileInfo().Mode()|os.ModeDir|100); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, header.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			_, err = io.Copy(file, tr)
			_ = file.Close()
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

		default:
			log.Debug("skip:", header)
		}
	}
	return nil
}

type StorageConfig struct {
	Type  string `mapstructure:"type" yaml:"type"`
	Local struct {
		// Dir 保存ディレクトリ
		Dir string `mapstructure:"dir" yaml:"dir"`
	} `mapstructure:"local" yaml:"local"`
	S3 struct {
		// Bucket バケット名
		Bucket       string `mapstructure:"bucket" yaml:"bucket"`
		Endpoint     string `mapstructure:"endpoint" yaml:"endpoint"`
		Region       string `mapstructure:"region" yaml:"region"`
		AccessKey    string `mapstructure:"accessKey" yaml:"accessKey"`
		AccessSecret string `mapstructure:"accessSecret" yaml:"accessSecret"`
	} `mapstructure:"s3" yaml:"s3"`
	Swift struct {
		// UserName ユーザー名
		UserName string `mapstructure:"username" yaml:"username"`
		// APIKey APIキー(パスワード)
		APIKey string `mapstructure:"apiKey" yaml:"apiKey"`
		// TenantName テナント名
		TenantName string `mapstructure:"tenantName" yaml:"tenantName"`
		// TenantID テナントID
		TenantID string `mapstructure:"tenantId" yaml:"tenantId"`
		// Container コンテナ名
		Container string `mapstructure:"container" yaml:"container"`
		// AuthURL 認証エンドポイント
		AuthURL string `mapstructure:"authUrl" yaml:"authUrl"`
	} `mapstructure:"swift" yaml:"swift"`
}
