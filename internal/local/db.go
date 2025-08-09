package local

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(databaseURL string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	return db
}
