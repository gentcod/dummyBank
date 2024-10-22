package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/gentcod/DummyBank/api"
	"github.com/gentcod/DummyBank/gapi"
	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/pb"
	"github.com/gentcod/DummyBank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"github.com/rakyll/statik/fs"
	_ "github.com/gentcod/DummyBank/doc/statik"

	// "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := util.LoadConfig("./app.env")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBUrl)
	if err != nil {
		log.Fatal("Couldn't connect to db:", err)
	}

	// runDBMigration(config.MigrationUrl, config.DBUrl)

	store := db.NewStore(conn)
	// runGinServer(config, store)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

// runGinServer initializes HTTP server.
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Couldn't initialize the server:", err)
	}

	err = server.Start(config.PortAddress)
	if err != nil {
		log.Fatal("Couldn't start up server:", err)
	}
}

// runGrpcServer initializes gRPC server.
func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Couldn't initialize the server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDummyBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcAddress)
	if err != nil {
		log.Fatal("Couldn't create listener:", err)
	}

	log.Printf("gRPC server is running on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Couldn't start gRPC server:", err)
	}
}

// runGatewayServer initializes gRPC HTTP gateway server.
func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Couldn't initialize the server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterDummyBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("Couldn't register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.PortAddress)
	if err != nil {
		log.Fatal("Couldn't create listener:", err)
	}

	log.Printf("gRPC HTTP gateway server is running on %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Couldn't start gRPC gateway server:", err)
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
