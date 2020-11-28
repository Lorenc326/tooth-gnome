package orm

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"strconv"
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

func (u *User) InsertIfNotExist(db *pg.DB) (bool, error) {
	return db.Model(u).SelectOrInsert()
}

func (u *User) SetReminders(db *pg.DB) (pg.Result, error) {
	return db.Model(u).
		Column("morning_time").
		Column("evening_time").
		WherePK().
		Update()
}

func (u *User) Train(db *pg.DB) (pg.Result, error) {
	return db.Model(u).
		Column("last_trained").
		WherePK().
		Update()
}

func (_ *User) GetUsersToRemind(db *pg.DB, users *[]User, now string, offset int, limit int) error {
	return db.Model(users).
		Column("id").
		Where("morning_time = ?", now).
		WhereOr("evening_time = ?", now).
		Offset(offset).
		Limit(limit).
		Select()
}

func (u *User) String() string {
	return fmt.Sprintf("User<id=%d lng=%s created_at=%s progress=%d>", u.ID, u.Lng, u.CreatedAt, u.Progress)
}

// Retrieve id for telegram interface
func (u User) Recipient() string {
	return strconv.Itoa(u.ID)
}
