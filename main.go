package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gentcod/DummyBank/api"
	_ "github.com/gentcod/DummyBank/doc/statik"
	"github.com/gentcod/DummyBank/gapi"
	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/pb"
	"github.com/gentcod/DummyBank/util"
	"github.com/gentcod/DummyBank/worker"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	// "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := util.LoadConfig("./app.env")
	if err != nil {
		log.Error().AnErr("cannot load config", err)
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBUrl)
	if err != nil {
		log.Error().AnErr("Couldn't connect to db:", err)
	}

	// runDBMigration(config.MigrationUrl, config.DBUrl)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store)
	// runGinServer(config, store, taskDistributor)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}

// runGinServer initializes HTTP server.
func runGinServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Error().AnErr("Couldn't initialize the server:", err)
	}

	err = server.Start(config.PortAddress)
	if err != nil {
		log.Error().AnErr("Couldn't start up server:", err)
	}
}

// runGrpcServer initializes gRPC server.
func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Error().AnErr("Couldn't initialize the server:", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterDummyBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcAddress)
	if err != nil {
		log.Error().AnErr("Couldn't create listener:", err)
	}

	log.Info().Msgf("gRPC server is running on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Error().AnErr("Couldn't start gRPC server:", err)
	}
}

// runGatewayServer initializes gRPC HTTP gateway server.
func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Error().AnErr("Couldn't initialize the server:", err)
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
		log.Error().AnErr("Couldn't register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Error().AnErr("cannot create statik fs", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.PortAddress)
	if err != nil {
		log.Error().AnErr("Couldn't create listener:", err)
	}

	log.Info().Msgf("gRPC HTTP gateway server is running on %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Error().AnErr("Couldn't start gRPC gateway server:", err)
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().AnErr("failed to start task processor:", err)
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
