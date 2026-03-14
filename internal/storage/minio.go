package storage

import (
	"context"
	"net/url"
	"path"
	"time"

	"github.com/mfb27/luban/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIO struct {
	Client        *minio.Client
	Bucket        string
	PublicBaseURL string
}

func NewMinIO(cfg config.MinIOConfig) (*MinIO, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	m := &MinIO{
		Client:        client,
		Bucket:        cfg.Bucket,
		PublicBaseURL: cfg.PublicBaseURL,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		// Create bucket for dev convenience.
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *MinIO) PublicURL(objectKey string) (string, error) {
	base, err := url.Parse(m.PublicBaseURL)
	if err != nil {
		return "", err
	}
	base.Path = path.Join(base.Path, m.Bucket, objectKey)
	return base.String(), nil
}

