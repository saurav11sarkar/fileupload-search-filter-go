package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func CreateToken(userID, role, tokenType, secret string, ttl time.Duration) (string, error) {
	claim := Claim{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "myapp",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Subject:   userID,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
}

func ParseToken(tokenString, secret, expectedType string) (*Claim, error) {
    token, err := jwt.ParseWithClaims(
        tokenString,
        &Claim{},
        func(t *jwt.Token) (any, error) {
            if t.Method != jwt.SigningMethodHS256 {
                return nil, fmt.Errorf(
                    "unexpected signing method: %v",
                    t.Header["alg"],
                )
            }
            return []byte(secret), nil
        },
    )

    if err != nil {
        return nil, err
    }

    claim, ok := token.Claims.(*Claim)
    if !ok || !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    if claim.Type != expectedType {
        return nil, fmt.Errorf("invalid token type")
    }

    return claim, nil
}
