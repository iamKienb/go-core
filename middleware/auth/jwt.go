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
}

func NewJWTGenerator(cfg configx.JwtConfig) Generator {
	privBytes, _ := os.ReadFile("private.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privBytes)

	pubBytes, _ := os.ReadFile("public.pem")
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pubBytes)

	return &JWTGenerator{
		privateKey: privKey,
		publicKey:  pubKey,
		cfg:        cfg,
	}
}

func (g *JWTGenerator) GeneratePair(claims Claims) (*Pair, error) {
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
