package gapi

import (
	"context"
	db "github.com/auronvila/simple-bank/db/sqlc"
	pb "github.com/auronvila/simple-bank/pb/user"
	"github.com/auronvila/simple-bank/util"
	"github.com/auronvila/simple-bank/val"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := ValidateCreateUser(req)
	if violations != nil {
		return nil, InvalidArgument(violations)
	}
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password could not be hashed!! %s", err)
	}
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		Email:          req.GetEmail(),
		FullName:       req.GetFullName(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user with these credentials already exist %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create the user!! %s", err)
	}

	rsp := &pb.CreateUserResponse{User: convertUser(user)}
	return rsp, nil
}

func ValidateCreateUser(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, FieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, FieldViolation("password", err))
	}

	if err := val.ValidateUserFullName(req.GetFullName()); err != nil {
		violations = append(violations, FieldViolation("full_name", err))
	}

	if err := val.ValidateEmailAddress(req.GetEmail()); err != nil {
		violations = append(violations, FieldViolation("email", err))
	}

	return violations
}
