package db

import (
	"database/sql"
	"log"
	"os"
	_ "simplebank/util"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@172.27.193.28/simple_bank?sslmode=disable"
)

// 在所有的单元测试中会用到
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	testQueries = New(testDB)
	//Run runs the tests. It returns an exit code to pass to os.Exit.
	m.Run()
	os.Exit(m.Run())
}
