package db

import (
	"context"
	"github.com/go-pg/pg/v10"
)

type tableInitiator interface {
	initiateTable(db *pg.DB) error
}

var models = []tableInitiator{
	(*User)(nil),
}

func ConnectDB(url string) *pg.DB {
	opt, err := pg.ParseURL(url)
	if err != nil {
		panic(err)
	}

	db := pg.Connect(opt)
	defer func() {
		// close connection and pass error next
		if err := recover(); err != nil {
			db.Close()
			panic(err)
		}
	}()

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	for _, model := range models {
		if err := model.initiateTable(db); err != nil {
			panic(err)
		}
	}

	return db
}
