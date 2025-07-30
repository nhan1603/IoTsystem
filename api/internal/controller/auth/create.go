package auth

import (
	"context"
	"fmt"

	"github.com/nhan1603/IoTsystem/api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (i impl) CreateUser(ctx context.Context, user *model.User) error {
	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// Create the user in the database
	err = i.repo.User().Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
