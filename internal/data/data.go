package data

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/data/ent"

	"github.com/go-kratos/kratos/v3/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewTodoRepo)

// Data holds the long-lived storage clients shared by repos.
type Data struct {
	db *ent.Client
}

// NewData opens the database client and returns it with a cleanup function.
func NewData(c *conf.Data) (*Data, func(), error) {
	dc := c.GetDatabase()
	db, err := ent.Open(dc.GetDriver(), dc.GetSource())
	if err != nil {
		return nil, nil, err
	}
	if dc.GetDebug() {
		db = db.Debug()
	}
	// Auto migration is a convenience for local development. In production,
	// apply schema changes as a separate reviewed step instead.
	if dc.GetAutoMigrate() {
		if err := db.Schema.Create(context.Background()); err != nil {
			db.Close()
			return nil, nil, err
		}
	}
	cleanup := func() {
		log.Info("closing the data resources")
		if err := db.Close(); err != nil {
			log.Error("failed closing the database", "err", err)
		}
	}
	return &Data{db: db}, cleanup, nil
}
