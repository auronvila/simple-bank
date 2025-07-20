package gapi

import (
	"context"
	db "github.com/auronvila/simple-bank/db/sqlc"
	pb "github.com/auronvila/simple-bank/pb/user"
	"github.com/auronvila/simple-bank/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := ValidateVerifyEmail(req)
	if violations != nil {
		return nil, InvalidArgument(violations)
	}

	txRes, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	rsp := &pb.VerifyEmailResponse{IsVerified: txRes.User.IsEmailVerified}
	return rsp, nil
}

func ValidateVerifyEmail(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, FieldViolation("email_id", err))
	}

	if err := validator.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, FieldViolation("secret_code", err))
	}

	return violations
}
