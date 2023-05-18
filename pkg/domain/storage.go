package domain

import (
	"io"
	"path/filepath"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/util/tarfs"
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
}

func buildLogPath(buildID string) string {
	const buildLogDirectory = "buildlogs"
	return filepath.Join(buildLogDirectory, buildID)
}

func SaveBuildLog(s Storage, buildID string, src io.Reader) error {
	err := s.Save(buildLogPath(buildID), src)
	if err != nil {
		return errors.Wrap(err, "saving build log")
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

func artifactPath(artifactID string) string {
	return filepath.Join("artifacts", artifactFilename(artifactID))
}

func artifactFilename(artifactID string) string {
	return artifactID + ".tar"
}

// SaveArtifact Artifactをtar形式で保存する
func SaveArtifact(s Storage, artifactID string, src io.Reader) error {
	if err := s.Save(artifactPath(artifactID), src); err != nil {
		return errors.Wrap(err, "failed to save artifact")
	}
	return nil
}

func GetArtifact(s Storage, artifactID string) (filename string, b []byte, err error) {
	r, err := s.Open(artifactPath(artifactID))
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to open artifact")
	}
	defer r.Close()
	b, err = io.ReadAll(r)
	return artifactFilename(artifactID), b, err
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

	return tarfs.Extract(inputFile, destPath)
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
