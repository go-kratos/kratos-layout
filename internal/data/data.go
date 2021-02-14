package data

import (
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Config is data config.
type Config struct {
	Driver string
	Source string
}

// Data .
type Data struct {
	// TODO warpped database client
}

// NewData .
func NewData(c *Config) (*Data, error) {
	return &Data{}, nil
}
