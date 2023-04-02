package domain

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"github.com/friendsofgo/errors"

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

func buildLogPath(buildID string) string {
	const buildLogDirectory = "buildlogs"
	return filepath.Join(buildLogDirectory, buildID)
}

func SaveBuildLog(s Storage, filename string, buildID string) error {
	if err := s.Move(filename, buildLogPath(buildID)); err != nil {
		return errors.Wrap(err, "failed to move build log")
	}
	return nil
}

func GetBuildLog(s Storage, buildID string) ([]byte, error) {
	r, err := s.Open(buildLogPath(buildID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open build log")
	}
	defer r.Close()
	return io.ReadAll(r)
}

func DeleteBuildLog(s Storage, buildID string) error {
	err := s.Delete(buildLogPath(buildID))
	if err != nil {
		return errors.Wrap(err, "failed to delete build log")
	}
	return nil
}

func artifactPath(id string) string {
	const artifactDirectory = "artifacts"
	return filepath.Join(artifactDirectory, id+".tar")
}

// SaveArtifact Artifactをtar形式で保存する
func SaveArtifact(s Storage, filename string, artifactID string) error {
	if err := s.Move(filename, artifactPath(artifactID)); err != nil {
		return errors.Wrap(err, "failed to save artifact")
	}
	return nil
}

func DeleteArtifact(s Storage, artifactID string) error {
	err := s.Delete(artifactPath(artifactID))
	if err != nil {
		return errors.Wrap(err, "failed to delete artifact")
	}
	return nil
}

// ExtractTarToDir tarファイルをディレクトリに展開する
func ExtractTarToDir(s Storage, artifactID string, destPath string) error {
	inputFile, err := s.Open(artifactPath(artifactID))
	if err != nil {
		return errors.Wrap(err, "couldn't open source file")
	}
	defer inputFile.Close()

	tr := tar.NewReader(inputFile)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return errors.Wrap(err, "bad tar file")
		}

		path := filepath.Join(destPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, header.FileInfo().Mode()); err != nil {
				return errors.Wrap(err, "failed to create directory")
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), header.FileInfo().Mode()|os.ModeDir|100); err != nil {
				return errors.Wrap(err, "failed to create directory")
			}

			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, header.FileInfo().Mode())
			if err != nil {
				return errors.Wrap(err, "failed to create file")
			}
			_, err = io.Copy(file, tr)
			_ = file.Close()
			if err != nil {
				return errors.Wrap(err, "failed to write file")
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
