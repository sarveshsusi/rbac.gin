package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TwoFATokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func Generate2FAToken(
	userID uuid.UUID,
	secret string,
) (string, error) {

	claims := TwoFATokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func Parse2FAToken(
	raw string,
	secret string,
) (*TwoFATokenClaims, error) {

	token, err := jwt.ParseWithClaims(
		raw,
		&TwoFATokenClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, err
	}

	return token.Claims.(*TwoFATokenClaims), nil
}
