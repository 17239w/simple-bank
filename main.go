package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@172.27.193.28/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	//func sql.Open(driverName string, dataSourceName string) (*sql.DB, error)
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
