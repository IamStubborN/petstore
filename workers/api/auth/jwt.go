package auth

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"path"
	"time"

	"github.com/IamStubborN/petstore/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/patrickmn/go-cache"

	"go.uber.org/zap"
)

type Claims struct {
	*jwt.StandardClaims
	Session
}

type Session struct {
	SessionID    string
	UserID       int64
	AllowMethods string
}

type JWTAuthService struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	blackList  *cache.Cache
	ttl        time.Duration
}

var jwtAuth JWTAuthService

func InitJWTAuth(cfg *config.Config) {
	generatePEMKeys(cfg.JWT.KeysPath)
	jwtAuth = JWTAuthService{
		publicKey:  loadPublicKey(cfg.JWT.KeysPath),
		privateKey: loadPrivateKey(cfg.JWT.KeysPath),
		blackList:  cache.New(cfg.JWT.TTL, cfg.JWT.TTL),
		ttl:        cfg.JWT.TTL,
	}
}

func AddToBlackList(token string) {
	jwtAuth.blackList.SetDefault(token, struct{}{})
}

func IsTokenInBlackList(token string) bool {
	_, isExist := jwtAuth.blackList.Get(token)

	return isExist
}

func loadPublicKey(pathToFolder string) *rsa.PublicKey {
	pk, err := ioutil.ReadFile(path.Join(pathToFolder, "public.pem"))
	if err != nil {
		zap.L().Fatal("can't read public.pem", zap.Error(err))
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pk)
	if err != nil {
		zap.L().Fatal("can't parse public.pem", zap.Error(err))
	}

	return publicKey
}

func loadPrivateKey(pathToFolder string) *rsa.PrivateKey {
	pk, err := ioutil.ReadFile(path.Join(pathToFolder, "private.pem"))
	if err != nil {
		zap.L().Fatal("can't read private.pem", zap.Error(err))
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pk)
	if err != nil {
		zap.L().Fatal("can't parse private.pem", zap.Error(err))
	}
	return privateKey
}

func GenerateToken(sessionID string, userID int64, allowMethods string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtAuth.ttl).Unix(),
		},
		Session: Session{
			SessionID:    sessionID,
			UserID:       userID,
			AllowMethods: allowMethods,
		},
	})

	return token.SignedString(jwtAuth.privateKey)
}

func ParseToken(jwtToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtAuth.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
