package biz

import (
	"context"

	v1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.
type Greeter struct {
	Hello string
}

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	CreateGreeter(context.Context, *Greeter) (*Greeter, error)
	UpdateGreeter(context.Context, *Greeter) (*Greeter, error)
}

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo GreeterRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, log: log.NewHelper(logger)}
}

// Create creates a Greeter, and returns the new Greeter.
func (uc *GreeterUsecase) Create(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Infof("Create: %v", g.Hello)
	return uc.repo.CreateGreeter(ctx, g)
}

// Update updates a Greeter.
func (uc *GreeterUsecase) Update(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Infof("Update: %v", g.Hello)
	return uc.repo.UpdateGreeter(ctx, g)
}
