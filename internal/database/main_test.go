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
	dbSrc = "postgres://postgres:oyetide50@localhost:5432/dummybank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSrc)
	if err != nil {
		log.Fatal("Couldn't connect to db:", err)
	}
	
	testQueries = New(conn)

	//Initialize connection test, terminate test if error occurs
	os.Exit(m.Run()) 
}