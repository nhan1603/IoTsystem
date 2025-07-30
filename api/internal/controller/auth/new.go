package auth

import (
	"context"

	"github.com/nhan1603/IoTsystem/api/internal/appconfig/iam"
	"github.com/nhan1603/IoTsystem/api/internal/model"
	"github.com/nhan1603/IoTsystem/api/internal/repository"
)

// Controller represents the specification of this pkg
type Controller interface {
	CheckAuth(ctx context.Context, inp LoginInput) (model.User, string, error)
	CreateUser(ctx context.Context, user *model.User) error
}

// New initializes a new Controller instance and returns it
func New(repo repository.Registry, config iam.Config) Controller {
	return impl{
		repo:   repo,
		config: config,
	}
}

type impl struct {
	repo   repository.Registry
	config iam.Config
}
