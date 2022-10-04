package models

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

var db *bun.DB

func CreateDBConnection(pgDB *sql.DB) {
	// Create a Bun db on top of it.
	db = bun.NewDB(pgDB, pgdialect.New())

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(false)))
}
