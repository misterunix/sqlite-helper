package helper

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

///
///
///

// open the sqlite database
func (db *DbConfig) OpenDB() error {

	var err error

	// create the directory if it does not exist
	os.MkdirAll(db.Path, os.ModePerm)

	//pragma := "?_pragma=journal_mode(WAL)&_pragma=read_uncommitted(true)"
	pragma := ""
	dbfile := db.Path + "/" + db.Filename + pragma

	db.Db, err = sql.Open("sqlite", dbfile)
	if err != nil {
		return err
	}

	return nil

}

// execute the sql statement, locking the function.
// locking is used to  prevent reading before the previous write is complete.
// this has a fallback timer on/for errors.
func (db *DbConfig) RunSQLStatement(sql string) error {

	// lock the function
	db.dblock.Lock()

	// unlock the function when leaving the function
	defer db.dblock.Unlock()

	// wait time for when there is an error or lock
	waitTime := 1

	statement, err := db.Db.Prepare(sql)
	if err != nil {
		return errors.New("Failed to prepare statement.")
	}

	defer statement.Close() // Close the statement when leaving the function.

	rounds := 0 // this is a hack to prevent an infinite loop

	for {

		if rounds > 10 {
			return errors.New("Failed to execute statement. rounds exceeded.")
		}

		_, err = statement.Exec()
		if err != nil {
			time.Sleep(time.Duration(waitTime) * time.Second) // backoff waiting for the db
			waitTime *= 2
			rounds++
		} else {
			break
		}
	}

	return nil
}

///
///
///

// creating the table from the struct, executing the sql statement
func (db *DbConfig) CreateTable(table string, s interface{}) error {

	sqlstatement := db.CreateTableFromStruct(table, s)

	err := db.RunSQLStatement(sqlstatement)
	if err != nil {
		return err
	}

	return nil
}

///
///
///

// creating the table from the struct, executing the sql statement
// dropping the table first
func (db *DbConfig) RemoveAndCreateNewDB(table string, s interface{}) error {

	sqlstatement := "DROP TABLE IF EXISTS " + table + ";"

	err := db.RunSQLStatement(sqlstatement)
	if err != nil {
		return err
	}

	err = db.CreateTable(table, s)
	if err != nil {
		return err
	}

	return nil
}

///
///
///

// build the sql statement to create a table from a struct
// does not drop if table already exsists
func (db *DbConfig) CreateTableFromStruct(table string, s interface{}) string {

	var reflectedValue reflect.Value = reflect.ValueOf(s) // reflect the struct (interface)

	var sqlstatement string // the sql statement to return

	sqlstatement1 := "CREATE TABLE IF NOT EXISTS " + table + " ("

	for i := 0; i < reflectedValue.NumField(); i++ {
		var vt string
		varName := reflectedValue.Type().Field(i).Name // get the name of the field
		sqlstatement += "," + varName + " "
		varType := reflectedValue.Type().Field(i).Type // get the type of the field

		// Did this differnt than the other reflect code. This is a work in progress.
		switch varType.Kind() {
		case reflect.Int:
			if varName == "ID" { // detect if the field is the ID field
				vt = "INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT"
			} else {
				vt = "INTEGER"
			}
		case reflect.Int8:
			vt = "INTEGER"
		case reflect.Int16:
			vt = "INTEGER"
		case reflect.Int32:
			vt = "INTEGER"
		case reflect.Int64:
			vt = "INTEGER"
		case reflect.Uint:
			vt = "INTEGER"
		case reflect.Uint8:
			vt = "INTEGER"
		case reflect.Uint16:
			vt = "INTEGER"
		case reflect.Uint32:
			vt = "INTEGER"
		case reflect.Uint64:
			vt = "INTEGER"
		case reflect.String:
			vt = "TEXT"
		case reflect.Float64:
			vt = "REAL"
		case reflect.Float32:
			vt = "REAL"
		case reflect.Bool:
			vt = "INTEGER"
		default:
			vt = "TEXT"
		}
		sqlstatement += vt
	}

	// such a crappy way to do this. Return to this at a later date.
	sqlstatement = sqlstatement[1:] // remove the first comma
	sqlstatement += ")"
	sqlstatement = sqlstatement1 + sqlstatement

	return sqlstatement
}

///
///
///

// Make insert into table from struct
func (db *DbConfig) InsertIntoTable(table string, s interface{}) string {

	var middlesql1 string
	var middlesql2 string

	var reflectedValue reflect.Value = reflect.ValueOf(s)

	middlesql1 = "INSERT INTO " + table + " ("
	middlesql2 = ")VALUES("
	for i := 0; i < reflectedValue.NumField(); i++ {

		varName := reflectedValue.Type().Field(i).Name
		varType := reflectedValue.Type().Field(i).Type
		varValue := reflectedValue.Field(i).Interface()

		if varName == "ID" {
			continue
		}

		middlesql1 += varName + ","

		// This is my normal way of working with reflect. Strings may be slower but easier to read.
		switch varType.Kind() {
		case reflect.Int:
			middlesql2 += fmt.Sprintf("%d", varValue.(int)) + ","
		case reflect.Int8:
			middlesql2 += fmt.Sprintf("%d", varValue.(int8)) + ","
		case reflect.Int16:
			middlesql2 += fmt.Sprintf("%d", varValue.(int16)) + ","
		case reflect.Int32:
			middlesql2 += fmt.Sprintf("%d", varValue.(int32)) + ","
		case reflect.Int64:
			middlesql2 += fmt.Sprintf("%d", varValue.(int64)) + ","
		case reflect.Uint:
			middlesql2 += fmt.Sprintf("%d", varValue.(uint)) + ","
		case reflect.Uint8:
			middlesql2 += fmt.Sprintf("%d", varValue.(uint8)) + ","
		case reflect.Uint16:
			middlesql2 += fmt.Sprintf("%d", varValue.(uint16)) + ","
		case reflect.Uint32:
			middlesql2 += fmt.Sprintf("%d", varValue.(uint32)) + ","
		case reflect.Uint64:
			middlesql2 += fmt.Sprintf("%d", varValue.(uint64)) + ","
		case reflect.String:
			middlesql2 += "'" + varValue.(string) + "',"
		case reflect.Float32:
			middlesql2 += fmt.Sprintf("%f", varValue.(float64)) + ","
		case reflect.Float64:
			middlesql2 += fmt.Sprintf("%f", varValue.(float64)) + ","
		case reflect.Bool:
			middlesql2 += fmt.Sprintf("%v", varValue.(bool)) + ","
		case reflect.Slice:
			middlesql2 += "'"
			for _, kk := range varValue.([]string) {
				middlesql2 += kk + " "
			}
			middlesql2 = strings.TrimRight(middlesql2, " ")
			middlesql2 += "',"
		default:
			return ""
		}
	}

	middlesql1 = middlesql1[:len(middlesql1)-1]
	middlesql2 = middlesql2[:len(middlesql2)-1] + ");"

	sqlreturn := middlesql1 + middlesql2

	return sqlreturn
}

///
///
///
