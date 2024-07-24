package helper

import (
	"database/sql"
	"sync"
)

type DbConfig struct {
	Db       *sql.DB    // db connection
	dblock   sync.Mutex // lock for the db.
	Path     string     // path to the db
	Filename string     // filename of the db
}

const VERSION = "0.0.1"
