package database

import (
	"database/sql"
)

type dbConn struct {
	conn *sql.DB
	error
}

//func New() SQLDatabase {
//	return &dbConn{}
//}
//
//func (db *dbConn) Open() {
//	connectionStr := os.Getenv("POSTGRESQL_CONNECTION_STRING")
//	driverName := os.Getenv("DRIVER_NAME")
//
//	db, err := sql.Open(driverName, connectionStr)
//	defer func(db *sql.DB) {
//		er := db.Close()
//		if er != nil {
//			panic(err.Error())
//		}
//	}(db)
//	if err != nil {
//		panic(err.Error())
//	}
//
//}
