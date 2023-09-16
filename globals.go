package helper

import (
	"database/sql"
	"sync"
)

type DbConfig struct {
	Db       *sql.DB
	dblock   sync.Mutex // lock for the db
	Path     string
	Filename string
}
