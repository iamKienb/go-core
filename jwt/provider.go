package jwtx

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
	configx "github.com/iamKienb/go-core/config"
)

type JWTX struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	cfg        configx.JwtConfig
}

func New(cfg configx.JwtConfig) (JWTXService, error) {
	generator := &JWTX{cfg: cfg}

	privBytes, err := os.ReadFile("private.pem")
	if err != nil {
		return nil, fmt.Errorf("jwt: read private key: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, fmt.Errorf("jwt: parse private key: %w", err)
	}

	pubBytes, err := os.ReadFile("public.pem")
	if err != nil {
		return nil, fmt.Errorf("jwt: read public key: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("jwt: parse public key: %w", err)
	}

	generator.privateKey = privKey
	generator.publicKey = pubKey

	return generator, nil
}
