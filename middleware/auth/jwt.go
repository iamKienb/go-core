package auth

import (
	"errors"
	"fmt"
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
	cfg configx.JwtConfig
}

func NewJWTGenerator(cfg configx.JwtConfig) TokenGenerator {
	return &JWTGenerator{
		cfg: cfg,
	}
}

func (g *JWTGenerator) GeneratePair(claims TokenClaims) (*TokenPair, error) {
	now := time.Now().UTC()
	accessExpAt := now.Add(g.cfg.AccessExpiry)
	refreshExpAt := now.Add(g.cfg.RefreshExpiry)

	accessToken, err := g.Sign(claims, accessExpAt, g.cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("jwt: sign access token: %w", err)
	}

	refreshToken, err := g.Sign(claims, refreshExpAt, g.cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("jwt: sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExpAt,
		RefreshTokenExpiresAt: refreshExpAt,
	}, nil
}

func (g *JWTGenerator) Sign(claims TokenClaims, expiryAt time.Time, secretKey string) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	claim := jwtClaims{
		UserID:          claims.UserId,
		Email:           claims.Email,
		Roles:           claims.Roles,
		PasswordVersion: claims.PasswordVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    "user-command",
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claim).SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (g *JWTGenerator) Verify(tokenString string) (*TokenClaims, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(g.cfg.PublicKey))
	if err != nil {
		return nil, fmt.Errorf("jwt: invalid public key: %w", err)
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("jwt: token is expired")
		}
		return nil, fmt.Errorf("jwt: parse failed: %w", err)
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return &TokenClaims{
			UserId:          claims.UserID,
			Email:           claims.Email,
			Roles:           claims.Roles,
			PasswordVersion: claims.PasswordVersion,
		}, nil
	}

	return nil, fmt.Errorf("jwt: invalid token")
}
