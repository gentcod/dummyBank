package gapi

import (
	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
		UpdatedAt:         timestamppb.New(user.UpdatedAt),
	}
}
