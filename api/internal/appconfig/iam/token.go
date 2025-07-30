package iam

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	crypto_helper "github.com/nhan1603/IoTsystem/api/internal/pkg/cryptos"
	pkgerrors "github.com/pkg/errors"
)

// JWTClaim represents jwt claim info
type JWTClaim struct {
	HostProfile HostProfile `json:"host_profile"`
	jwt.StandardClaims
}

// GenerateToken generate tokens
func (cfg Config) GenerateToken(claim JWTClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(cfg.JWTKey))
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	encryptionKey := os.Getenv("CYPHER_KEY")
	encryptedToken, err := crypto_helper.EncryptMessage([]byte(encryptionKey), tokenString)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	return encryptedToken, nil
}

// ValidateToken validates token
func (cfg Config) ValidateToken(entryptedToken string) (*jwt.Token, error) {
	encryptionKey := os.Getenv("CYPHER_KEY")
	strToken, err := crypto_helper.DecryptMessage([]byte(encryptionKey), entryptedToken)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	token, err := jwt.ParseWithClaims(
		strToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTKey), nil
		},
	)

	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token expired")
	}

	return token, nil
}

// GetUserProfileFomToken extracts user profile from token
func GetUserProfileFomToken(token *jwt.Token) (HostProfile, error) {
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return HostProfile{}, errors.New("token claims missing")
	}

	return claims.HostProfile, nil
}
