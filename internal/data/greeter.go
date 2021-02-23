package data

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos/kratos-layout/internal/biz"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper("data/greeter", logger),
	}
}

func (r *greeterRepo) CreateGreeter(g *biz.Greeter) error {
	return nil
}

func (r *greeterRepo) UpdateGreeter(g *biz.Greeter) error {
	return nil
}
