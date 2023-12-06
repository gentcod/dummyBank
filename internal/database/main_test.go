package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSrc = "postgres://root:secret@localhost:5432/dummy_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSrc)
	if err != nil {
		log.Fatal("Couldn't connect to db:", err)
	}
	
	testQueries = New(testDB)

	//Initialize connection test, terminate test if error occurs
	os.Exit(m.Run()) 
}