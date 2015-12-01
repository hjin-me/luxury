package db

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbIns               *sql.DB
	createErr           error
	ErrConnectNotCreate = errors.New("connection not create")
)

func Create(dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		createErr = err
		return err
	}

	dbIns = db
	return nil
}

func Close() error {
	if dbIns != nil {
		return dbIns.Close()
	}
	return createErr
}

func Connect() (*sql.DB, error) {
	if dbIns == nil {
		return nil, ErrConnectNotCreate
	}
	return dbIns, nil
}
