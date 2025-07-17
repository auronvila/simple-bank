package gapi

import (
	"context"
	db "github.com/auronvila/simple-bank/db/sqlc"
	pb "github.com/auronvila/simple-bank/pb/user"
	"github.com/auronvila/simple-bank/util"
	"github.com/auronvila/simple-bank/validator"
	"github.com/auronvila/simple-bank/worker"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
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
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			Email:          req.GetEmail(),
			FullName:       req.GetFullName(),
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{Username: user.Username}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(3 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user with these credentials already exist %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create the user!! %s", err)
	}

	rsp := &pb.CreateUserResponse{User: convertUser(txResult.User)}
	return rsp, nil
}

func ValidateCreateUser(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, FieldViolation("username", err))
	}

	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, FieldViolation("password", err))
	}

	if err := validator.ValidateUserFullName(req.GetFullName()); err != nil {
		violations = append(violations, FieldViolation("full_name", err))
	}

	if err := validator.ValidateEmailAddress(req.GetEmail()); err != nil {
		violations = append(violations, FieldViolation("email", err))
	}

	return violations
}
