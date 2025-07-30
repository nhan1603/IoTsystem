package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nhan1603/IoTsystem/api/internal/appconfig/iam"
	"github.com/nhan1603/IoTsystem/api/internal/model"
	"github.com/nhan1603/IoTsystem/api/internal/repository/user"
	pkgerrors "github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// LoginInput represents input to authenticate
type LoginInput struct {
	Email, Password string
}

var (
	checkAuthFn = checkAuth
)

// CheckAuth handles authentication checking
func (i impl) CheckAuth(ctx context.Context, inp LoginInput) (model.User, string, error) {
	log.Printf("[CheckAuth] START checking authentication (email: %s)\n", inp.Email)
	u, err := checkAuthFn(ctx, i, inp)
	if err != nil {
		return model.User{}, "", err
	}
	log.Println("[CheckAuth] login successfully")

	token, err := generateJWTToken(i, u)
	if err != nil {
		return model.User{}, "", err
	}

	return u, token, nil
}

func checkAuth(ctx context.Context, i impl, inp LoginInput) (model.User, error) {
	u, err := i.repo.User().GetByEmail(ctx, user.GetUserInput{
		Email: inp.Email,
	})
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inp.Password)); err != nil {
		return model.User{}, ErrUserNotFound
	}

	return u, nil
}

func generateJWTToken(i impl, u model.User) (string, error) {
	log.Println("[CheckAuth] generating token...")

	token, err := i.config.GenerateToken(iam.JWTClaim{
		HostProfile: iam.HostProfile{
			ID:    u.ID,
			Email: u.Email,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(60 * time.Minute).Unix(),
		},
	})
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	log.Println("[CheckAuth] token generated successfully")

	return token, nil
}
