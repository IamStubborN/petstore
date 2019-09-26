package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path"

	"go.uber.org/zap"
)

func generatePEMKeys(pathToFolder string) {
	createFoldersIfNotExist(pathToFolder)

	p := path.Clean(pathToFolder)
	publicPath := path.Join(p, "public.pem")
	privatePath := path.Join(p, "private.pem")

	if isExistsFile(publicPath) && isExistsFile(privatePath) {
		return
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		zap.L().Fatal("can't generate key", zap.Error(err))
	}

	savePEMKey(privatePath, key)
	savePublicPEMKey(publicPath, key.PublicKey)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	if err != nil {
		zap.L().Fatal("can't create file", zap.Error(err))
	}
	defer deferError(outFile.Close)

	bytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		zap.L().Fatal("can't encode private pem key", zap.Error(err))
	}

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: bytes,
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		zap.L().Fatal("can't encode private pem key", zap.Error(err))
	}
}

func savePublicPEMKey(fileName string, publicKey rsa.PublicKey) {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		zap.L().Fatal("can't marshal public pem key", zap.Error(err))
	}

	var pemKey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemFile, err := os.Create(fileName)
	if err != nil {
		zap.L().Fatal("can't create file", zap.Error(err))
	}
	defer deferError(pemFile.Close)

	err = pem.Encode(pemFile, pemKey)
	if err != nil {
		zap.L().Fatal("can't encode public pem key", zap.Error(err))
	}
}

func deferError(f func() error) {
	if err := f(); err != nil {
		zap.L().Error("error in defer", zap.Error(err))
	}
}

func createFoldersIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0777); err != nil {
			zap.L().Info("mkdir error", zap.Error(err))
		}
	}
}

func isExistsFile(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
