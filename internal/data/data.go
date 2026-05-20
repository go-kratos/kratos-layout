package data

import (
	"github.com/go-kratos/kratos-layout/internal/conf"

	"github.com/go-kratos/kratos/v3/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewTodoRepo)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}
