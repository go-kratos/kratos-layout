package data

import (
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// DBConfig is data config.
type DBConfig struct {
	Driver string
	Source string
}

// Data .
type Data struct {
	// TODO warpped database client
}

// NewData .
func NewData(c *DBConfig) (*Data, error) {
	return &Data{}, nil
}
