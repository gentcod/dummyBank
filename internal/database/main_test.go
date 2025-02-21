package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/gentcod/DummyBank/util"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var testQueries *Queries
var testDB *sql.DB
var testRDB *redis.Client

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../app.env")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBUrl)
	if err != nil {
		log.Fatal("Couldn't connect to db:", err)
	}
	
	testQueries = New(testDB)

	testRDB = redis.NewClient(&redis.Options{
		Addr: config.RedisAddress,
	})

	//Initialize connection test, terminate test if error occurs
	os.Exit(m.Run()) 
}