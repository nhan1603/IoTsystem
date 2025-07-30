package user

import (
	"context"

	"github.com/nhan1603/IoTsystem/api/internal/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// Repository provides the specification of the functionality provided by this pkg
type Repository interface {
	GetByEmail(context.Context, GetUserInput) (model.User, error)
	Create(ctx context.Context, user *model.User) error
}

// New returns an implementation instance satisfying Repository
func New(dbConn boil.ContextExecutor) Repository {
	return impl{
		dbConn: dbConn,
	}

}

type impl struct {
	dbConn boil.ContextExecutor
}
