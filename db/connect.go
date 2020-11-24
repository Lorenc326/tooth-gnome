package db

import (
	"context"
	"github.com/go-pg/pg/v10"
)

func ConnectDB (url string) (*pg.DB, error) {
	opt, err := pg.ParseURL(url)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opt)
	defer db.Close()

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return db, nil
}