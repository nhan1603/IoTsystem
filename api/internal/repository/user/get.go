package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/nhan1603/IoTsystem/api/internal/model"
	"github.com/nhan1603/IoTsystem/api/internal/repository/dbmodel"
	pkgerrors "github.com/pkg/errors"
)

// GetUserInput represents for input to get user
type GetUserInput struct {
	Email string
}

// GetByCriteria retrieves user by input
func (i impl) GetByEmail(ctx context.Context, inp GetUserInput) (model.User, error) {
	o, err := dbmodel.Users(
		dbmodel.UserWhere.Email.EQ(strings.ToLower(inp.Email))).One(ctx, i.dbConn)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrNotFound
		}

		return model.User{}, pkgerrors.WithStack(err)
	}

	return toModelUser(o), nil
}
