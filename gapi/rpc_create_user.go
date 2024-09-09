package gapi

import (
	"context"

	"github.com/gentcod/DummyBank/pb"
	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/util"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	arg := db.CreateUserParams{
		ID: uuid.New(),
		Username: req.GetUsername(),
		FullName: req.GetFullName(),
		Email: req.GetEmail(),
		HarshedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name(){
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username/email is already taken: %v", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err) 
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}