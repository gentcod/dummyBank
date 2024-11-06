package gapi

import (
	"fmt"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/pb"
	"github.com/gentcod/DummyBank/token"
	"github.com/gentcod/DummyBank/util"
	"github.com/gentcod/DummyBank/worker"
)

//Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedDummyBankServer
	config util.Config
	store db.Store
	tokenGenerator token.Generator
	taskDistributor worker.TaskDistributor
}

//NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenGenerator, err := token.NewPasetoGenerator(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize token generator: %v", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenGenerator: tokenGenerator,
		taskDistributor: taskDistributor,
	}

	return server, nil
}