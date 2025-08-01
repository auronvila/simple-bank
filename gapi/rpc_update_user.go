package gapi

import (
	"context"
	"database/sql"
	"errors"
	db "github.com/auronvila/simple-bank/db/sqlc"
	pb "github.com/auronvila/simple-bank/pb/user"
	"github.com/auronvila/simple-bank/util"
	"github.com/auronvila/simple-bank/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update the data of another user")
	}

	violations := ValidateUpdateUser(req)
	if violations != nil {
		return nil, InvalidArgument(violations)
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "password could not be hashed!! %s", err)
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  req.Password != nil,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update the user!! %s", err)
	}

	rsp := &pb.UpdateUserResponse{User: convertUser(user)}
	return rsp, nil
}

func ValidateUpdateUser(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, FieldViolation("username", err))
	}

	// optional fields below
	if req.Password != nil {
		if err := validator.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, FieldViolation("password", err))
		}
	}

	if req.FullName != nil {
		if err := validator.ValidateUserFullName(req.GetFullName()); err != nil {
			violations = append(violations, FieldViolation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := validator.ValidateEmailAddress(req.GetEmail()); err != nil {
			violations = append(violations, FieldViolation("email", err))
		}
	}

	return violations
}
