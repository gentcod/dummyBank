package gapi

import (
	"context"
	"database/sql"

	"github.com/gentcod/DummyBank/extensions"
	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/pb"
	"github.com/gentcod/DummyBank/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, valiateParameters(violations)
	}
	
	user, err := server.store.GetUserWithPassword(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	
	err = util.CheckPassword(req.GetPassword(), user.HarshedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user authentication error: %v", err)
	}

	accessToken, accessPayload, err := server.tokenGenerator.CreateToken(user.Username, user.ID, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v", err)
	}

	refreshToken, refreshPayload, err := server.tokenGenerator.CreateToken(user.Username, user.ID, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %v", err)
	}

	metadata := server.extractMetadata(ctx)

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID: refreshPayload.ID,
		Username: user.Username,
		RefreshToken: refreshToken,
		UserAgent: metadata.UserAgent,
		ClientIp: metadata.ClientIP,
		IsBlocked: false,
		ExpiresAt: refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user session: %v", err)
	}

	rsp := &pb.LoginUserResponse{
		SessionId: session.ID.String(),
		AccessToken: accessToken,
		AccessTokenExpiredAt: timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken: refreshToken,
		RefreshTokenExpiredAt: timestamppb.New(refreshPayload.ExpiredAt),
		User: convertUser(user),
	}
	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := extensions.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := extensions.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return
}