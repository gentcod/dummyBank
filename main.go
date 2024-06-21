package main

import (
	"database/sql"
	"log"

	"github.com/gentcod/DummyBank/api"
	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/util"
	_ "github.com/lib/pq"

	// "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBUrl)
	if err != nil {
		log.Fatal("Couldn't connect to db:", err)
	}

	// runDBMigration(config.MigrationUrl, config.DBUrl)

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Couldn't initialize the server:", err)
	}

	err = server.Start(config.Port)
	if err != nil {
		log.Fatal("Couldn't start up server:", err)
	}
}

// runDBMigration runs DB migrations when building docker images
// func runDBMigration(migrationURL string, dbURL string) {
// 	migration, err := migrate.New(migrationURL, dbURL)
// 	if err != nil {
// 		log.Fatal("Failed to create migration instance", err)
// 	}

// 	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
// 		log.Fatal("Database migration failed", err)
// 	}

// 	log.Println("db migration successful")
// }
