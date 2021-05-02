package common

import (
	"fmt"
	"strings"

	storage2 "github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
)

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

func (c *StorageConfig) InitStorage() (storage2.Storage, error) {
	switch strings.ToLower(c.Type) {
	case "local":
		return storage2.NewLocalStorage(c.Local.Dir)
	case "s3":
		return storage2.NewS3Storage(c.S3.Bucket, c.S3.AccessKey, c.S3.AccessSecret, c.S3.Region, c.S3.Endpoint)
	case "swift":
		return storage2.NewSwiftStorage(c.Swift.Container, c.Swift.UserName, c.Swift.APIKey, c.Swift.TenantName, c.Swift.TenantID, c.Swift.AuthURL)
	default:
		return nil, fmt.Errorf("unknown storage: %s", c.Type)
	}
}
