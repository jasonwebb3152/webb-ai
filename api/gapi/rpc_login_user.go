package gapi

import (
	"context"
	"time"

	"github.com/jasonwebb3152/webb-ai/api/pb"
	db "github.com/jasonwebb3152/webb-ai/db/sqlc"
	"github.com/jasonwebb3152/webb-ai/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/** Using simplified, custom authorization with refresh and auth tokens */
func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	// verify fields?, get user, verify pass, make refresh, make access, persist session, return
	// Check username
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to login user")
	}

	// Verify password is correct
	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to login user")
	}

	// Create refresh token
	refreshToken, refreshPayload, err := server.maker.CreateToken(user.Username, time.Hour*24)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create refresh token: %s", err)
	}

	// Persist Session
	metadata := server.extractMetadata(ctx)
	sessionParams := db.CreateSessionParams{
		ID:           refreshPayload.Id,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIp,
		ExpiresAt:    refreshPayload.ExpiresAt,
	}
	session, err := server.store.CreateSession(ctx, sessionParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create session record: %s", err)
	}

	// Create Auth Token
	accessToken, accessPayload, err := server.maker.CreateToken(user.Username, time.Minute*30)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create auth token: %s", err)
	}

	return &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		RefreshToken:          session.RefreshToken,
		AccessToken:           accessToken,
		RefreshTokenExpiresAt: timestamppb.New(session.ExpiresAt),
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt),
	}, nil
}
