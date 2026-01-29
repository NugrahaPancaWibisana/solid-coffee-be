package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/apperror"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	*dto.JWTClaims
}

func NewJWTClaims(id int, role string) *JwtClaims {
	return &JwtClaims{
		JWTClaims: &dto.JWTClaims{
			UserID: id,
			Role:   role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
				Issuer:    os.Getenv("JWT_ISSUER"),
			},
		},
	}
}

func (jc *JwtClaims) GenToken() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", apperror.ErrSecretNotFound
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jc)
	return token.SignedString([]byte(jwtSecret))
}

func (jc *JwtClaims) VerifyToken(token string) (bool, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return false, apperror.ErrSecretNotFound
	}

	jwtToken, err := jwt.ParseWithClaims(token, jc, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, apperror.ErrTokenExpired
		}
		return false, err
	}

	if !jwtToken.Valid {
		return false, apperror.ErrTokenInvalid
	}

	iss, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		return false, apperror.ErrTokenClaimsInvalid
	}

	expectedIssuer := os.Getenv("JWT_ISSUER")
	if expectedIssuer == "" {
		return false, apperror.ErrIssuerNotFound
	}

	if iss != expectedIssuer {
		return false, apperror.ErrInvalidIssuer
	}

	return true, nil
}
