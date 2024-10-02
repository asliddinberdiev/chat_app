package db

import (
	"database/sql"
	"fmt"

	"github.com/asliddinberdiev/chat_app/conf"
	_ "github.com/lib/pq"
)

type Databse struct {
	db *sql.DB
}

func NewDatabse(conf conf.Postgres) (*Databse, error) {
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Database, conf.SSLMode)

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	return &Databse{db: db}, nil
}

func (d *Databse) Close() {
	d.db.Close()
}

func (d *Databse) GetDB() *sql.DB {
	return d.db
}
