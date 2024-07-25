package helper

import (
	"database/sql"
	"sync"
)

// DbConfig is a struct that holds the database connection, path & filename
// to the database.
type DbConfig struct {
	Db       *sql.DB    // db connection
	dblock   sync.Mutex // lock for the db.
	Path     string     // path to the db
	Filename string     // filename of the db
}

const VERSION = "0.0.1"
