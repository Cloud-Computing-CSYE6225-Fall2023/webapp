package database

import (
	"database/sql"
	"os"
)

type dbConn struct {
}

func New() SQLDatabase {
	return &dbConn{}
}

func (dbo *dbConn) Open() (*sql.DB, error) {
	connectionStr := os.Getenv("POSTGRESQL_CONNECTION_STRING")
	driverName := os.Getenv("DRIVER_NAME")

	db, err := sql.Open(driverName, connectionStr)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	return db, nil
}
