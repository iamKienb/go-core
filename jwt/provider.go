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

	if err := generator.loadPrivateKey(keyPath(cfg.PrivateKeyPath, "private.pem")); err != nil {
		return nil, err
	}
	if err := generator.loadPublicKey(keyPath(cfg.PublicKeyPath, "public.pem")); err != nil {
		return nil, err
	}

	return generator, nil
}

func NewVerifier(cfg configx.JwtConfig) (JWTXService, error) {
	verifier := &JWTX{cfg: cfg}
	if err := verifier.loadPublicKey(keyPath(cfg.PublicKeyPath, "public.pem")); err != nil {
		return nil, err
	}

	return verifier, nil
}

func (x *JWTX) loadPrivateKey(path string) error {
	privBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("jwt: read private key %q: %w", path, err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return fmt.Errorf("jwt: parse private key %q: %w", path, err)
	}

	x.privateKey = privKey
	return nil
}

func (x *JWTX) loadPublicKey(path string) error {
	pubBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("jwt: read public key %q: %w", path, err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return fmt.Errorf("jwt: parse public key %q: %w", path, err)
	}

	x.publicKey = pubKey
	return nil
}

func keyPath(path string, fallback string) string {
	if path != "" {
		return path
	}

	return fallback
}
