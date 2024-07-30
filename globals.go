package helper

import (
	"database/sql"
	"sync"
)

const VERSION = "0.0.2"

type DbConfig struct {
	Db       *sql.DB
	dblock   sync.Mutex // lock for the db
	Path     string
	Filename string
}
