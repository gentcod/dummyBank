package gapi

import (
	"context"
	"fmt"

	"github.com/gentcod/DummyBank/extensions"
	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmalRequest(req)
	if violations != nil {
		return nil, valiateParameters(violations)
	}

	res, err := server.store.VerfiyEmailTx(ctx, db.VerfiyEmailParams{
		ID:    req.GetId(),
		Token: req.GetToken(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify user email: %v", err)
	}

	rsp := &pb.VerifyEmailResponse{
		IsVerified: res.IsEmailVerified,
	}
	return rsp, nil
}

func validateVerifyEmalRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	fmt.Println(req.GetId(), req.GetToken())
	if err := extensions.ValidateUsername(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	if err := extensions.ValidateSecretCode(req.GetToken()); err != nil {
		violations = append(violations, fieldViolation("token", err))
	}

	return
}
