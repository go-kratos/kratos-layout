package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (u userRepo) Find(ctx context.Context, id int64) (*biz.User, error) {
	return nil, fmt.Errorf("implement me")
}

func (u userRepo) Save(ctx context.Context, user *biz.User) error {
	return fmt.Errorf("implement me")
}

func (u userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	return nil, fmt.Errorf("implement me")
}
