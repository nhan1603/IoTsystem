package user

import (
	"github.com/nhan1603/IoTsystem/api/internal/model"
	"github.com/nhan1603/IoTsystem/api/internal/repository/dbmodel"
)

func toModelUsers(orms dbmodel.UserSlice) []model.User {
	users := make([]model.User, len(orms))
	for i, o := range orms {
		users[i] = toModelUser(o)
	}
	return users
}

func toModelUser(user *dbmodel.User) model.User {
	return model.User{
		ID:       int64(user.ID),
		Email:    user.Email,
		Password: user.PasswordHash,
	}
}
