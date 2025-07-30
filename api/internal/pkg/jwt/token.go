package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	pkgerrors "github.com/pkg/errors"
)

// JWTClaim represents jwt claim info
type JWTClaim struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// JWTKey is key to generate JWT token
// TODO: add to env
var JWTKey = []byte("CCO3rTNkzj5ll9JogQW5")

// GenerateToken generate tokens
func GenerateToken(claim JWTClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	return tokenString, nil
}

// ValidateToken validates token
func ValidateToken(strToken string) error {
	token, err := jwt.ParseWithClaims(
		strToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return JWTKey, nil
		},
	)

	if err != nil {
		return pkgerrors.WithStack(err)
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return errors.New("token expired")
	}

	return nil
}
