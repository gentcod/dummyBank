package gapi

import (
	"fmt"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/pb"
	"github.com/gentcod/DummyBank/token"
	"github.com/gentcod/DummyBank/util"
)

//Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedDummyBankServer
	config util.Config
	store db.Store
	tokenGenerator token.Generator
}

//NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenGenerator, err := token.NewPasetoGenerator(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize token generator: %v", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenGenerator: tokenGenerator,
	}

	return server, nil
}