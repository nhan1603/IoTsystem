package user

import (
	"context"
	"log"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/nhan1603/IoTsystem/api/internal/model"
	"github.com/nhan1603/IoTsystem/api/internal/repository/dbmodel"
	pkgerrors "github.com/pkg/errors"
)

// Create create user by input
func (i impl) Create(ctx context.Context, user *model.User) error {
	dbUser := dbmodel.User{
		Email:        user.Email,
		PasswordHash: user.Password,
		Username:     user.Name,
	}
	err := dbUser.Insert(ctx, i.dbConn, boil.Infer())
	log.Println("Finish insert with values " + dbUser.Email)

	if err != nil {
		pkgerrors.WithStack(err)
	}

	return nil
}
