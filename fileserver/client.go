package fileserver

import (
	"strings"

	"github.com/IamStubborN/petstore/config"
	"github.com/minio/minio-go/v6"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type FileManager interface {
	PutFile(bucket, contentType, filePath string) error
}

type fm struct {
	client *minio.Client
}

var manager *fm

func InitMinio(cfg *config.Config) {
	minioClient, err := minio.New(
		cfg.FileServer.Endpoint+":"+cfg.FileServer.Port,
		cfg.FileServer.AccessKey,
		cfg.FileServer.SecretKey,
		cfg.FileServer.SSL)
	if err != nil {
		zap.L().Fatal("can't create clien to minio server", zap.Error(err))
	}

	manager = &fm{minioClient}
}

func GetFM() FileManager {
	return manager
}

func (fm fm) PutFile(bucket, contentType, filePath string) error {
	isBucketExist, err := fm.client.BucketExists(bucket)
	if err != nil {
		return errors.Wrap(err, "can't check bucket exist")
	}

	if !isBucketExist {
		if err = fm.client.MakeBucket(bucket, ""); err != nil {
			return errors.Wrap(err, "can't create bucket")
		}
	}

	opts := minio.PutObjectOptions{ContentType: contentType}

	res := strings.Split(filePath, "/")
	fileName := res[len(res)-1]

	n, err := fm.client.FPutObject(bucket, fileName, filePath, opts)
	if err != nil {
		return err
	}

	zap.L().Info("file uploaded",
		zap.String("file", fileName),
		zap.Int64("written bytes", n))

	return nil
}
