package gapi

import (
	"context"

	"github.com/jasonwebb3152/webb-ai/api/pb"
	db "github.com/jasonwebb3152/webb-ai/db/sqlc"
	"github.com/jasonwebb3152/webb-ai/util"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if violations := validateCreateUserRequest(req); len(violations) > 0 {
		return nil, InvalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		// Error handling comes from gRPC packages
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	dbReq := db.CreateUserParams{
		Username:       req.GetUsername(),
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
		HashedPassword: hashedPassword,
	}
	user, err := server.store.CreateUser(ctx, dbReq)
	if err != nil {
		errCode := db.ErrorCode(err)
		switch errCode {
		case db.UniqueViolation:
			return nil, status.Errorf(codes.AlreadyExists, "unique field already exists: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	return &pb.CreateUserResponse{
		User: convertUser(user),
	}, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := util.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := util.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := util.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	return
}
