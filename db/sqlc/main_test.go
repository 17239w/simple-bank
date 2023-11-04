package db

import (
	"database/sql"
	"log"
	"os"
	"simplebank/util"
	_ "simplebank/util"
	"testing"

	_ "github.com/lib/pq"
)

// 在所有的单元测试中会用到
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	testQueries = New(testDB)
	//Run runs the tests. It returns an exit code to pass to os.Exit.
	m.Run()
	os.Exit(m.Run())
}
