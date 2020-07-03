package postgres

import (
	"github.com/go-pg/pg/v9"
)

/*
Connect opens a connection to the postgres database
*/
func Connect(pgOptions pg.Options) *pg.DB {
	// open a connection to the database
	return pg.Connect(&pgOptions)
}
