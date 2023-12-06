package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const dbDriver = "postgres"

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	godotenv.Load("../../.env")
	dbSrc := os.Getenv("POSTGRES_DOCKER_CI_DB_URL")

	testDB, err = sql.Open(dbDriver, dbSrc)
	if err != nil {
		log.Fatal("Couldn't connect to db:", err)
	}
	
	testQueries = New(testDB)

	//Initialize connection test, terminate test if error occurs
	os.Exit(m.Run()) 
}