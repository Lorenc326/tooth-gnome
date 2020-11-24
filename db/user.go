package db

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type User struct {
	ID          int
	Lng         string
	MorningTime string
	EveningTime string
	LastTrained string
	CreatedAt   string
	Progress    int16
}

func (u *User) initiateTable(db *pg.DB) error {
	return db.Model(u).CreateTable(&orm.CreateTableOptions{IfNotExists: true})
}

func (u *User) String() string {
	return fmt.Sprintf("User<id=%d lng=%s created_at=%s progress=%d>", u.ID, u.Lng, u.CreatedAt, u.Progress)
}
