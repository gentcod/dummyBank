package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gentcod/DummyBank/api"
	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const dbDriver = "postgres"

func main() {
	godotenv.Load("prod.env")
	dbSrc := os.Getenv("POSTGRES_DOCKER_DB_URL")
	port := os.Getenv("PORT")
	portAddress := fmt.Sprintf("localhost:%v", port)

	conn, err := sql.Open(dbDriver, dbSrc)
	if err != nil {
		log.Fatal("Couldn't connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(portAddress)
	if err != nil {
		log.Fatal("Couldn't start up server:", err)
	}
}