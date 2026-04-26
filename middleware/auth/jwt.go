package authx

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	configx "github.com/iamKienb/shopify-go-platform/config"
)

type jwtClaims struct {
	UserID          string   `json:"uid"`
	Email           string   `json:"email"`
	Roles           []string `json:"roles"`
	PasswordVersion int      `json:"pwd_v"`
	jwt.RegisteredClaims
}

type JWTGenerator struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	cfg        configx.JwtConfig
	initErr    error
}

func NewJWTGenerator(cfg configx.JwtConfig) Generator {
	generator := &JWTGenerator{cfg: cfg}

	privBytes, err := os.ReadFile("private.pem")
	if err != nil {
		generator.initErr = fmt.Errorf("jwt: read private key: %w", err)
		return generator
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		generator.initErr = fmt.Errorf("jwt: parse private key: %w", err)
		return generator
	}

	pubBytes, err := os.ReadFile("public.pem")
	if err != nil {
		generator.initErr = fmt.Errorf("jwt: read public key: %w", err)
		return generator
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		generator.initErr = fmt.Errorf("jwt: parse public key: %w", err)
		return generator
	}

	generator.privateKey = privKey
	generator.publicKey = pubKey

	return generator
}

func (g *JWTGenerator) GeneratePair(claims Claims) (*Pair, error) {
	if g.initErr != nil {
		return nil, g.initErr
	}

	now := time.Now().UTC()
	accessExpAt := now.Add(g.cfg.AccessExpiry)
	refreshExpAt := now.Add(g.cfg.RefreshExpiry)

	accessToken, err := g.Sign(claims, accessExpAt)
	if err != nil {
		return nil, fmt.Errorf("jwt: sign access token: %w", err)
	}

	refreshToken, err := g.Sign(claims, refreshExpAt)
	if err != nil {
		return nil, fmt.Errorf("jwt: sign refresh token: %w", err)
	}

	return &Pair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExpAt,
		RefreshTokenExpiresAt: refreshExpAt,
	}, nil
}

func (g *JWTGenerator) Sign(claims Claims, expiryAt time.Time) (string, error) {
	if g.initErr != nil {
		return "", g.initErr
	}

	claim := jwtClaims{
		UserID:          claims.UserID,
		Email:           claims.Email,
		Roles:           claims.Roles,
		PasswordVersion: claims.PasswordVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    "user-command",
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claim).SignedString(g.privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (g *JWTGenerator) Verify(tokenString string) (*Claims, error) {
	if g.initErr != nil {
		return nil, g.initErr
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method: %v", t.Header["alg"])
		}
		return g.publicKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("jwt: token is expired")
		}
		return nil, fmt.Errorf("jwt: parse failed: %w", err)
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return &Claims{
			UserID:          claims.UserID,
			Email:           claims.Email,
			Roles:           claims.Roles,
			PasswordVersion: claims.PasswordVersion,
		}, nil
	}

	return nil, fmt.Errorf("jwt: invalid token")
}
