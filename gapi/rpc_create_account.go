package gapi

import (
	"context"
	db "github.com/auronvila/simple-bank/db/sqlc"
	pb "github.com/auronvila/simple-bank/pb/account"
	"github.com/auronvila/simple-bank/val"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := ValidateCreateAccount(req)
	if violations != nil {
		return nil, InvalidArgument(violations)
	}

	if authPayload.Username != req.GetOwner() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update the data of another user")
	}

	arg := db.CreateAccountParams{
		Owner:    req.GetOwner(),
		Currency: req.GetCurrency(),
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				return nil, status.Error(codes.AlreadyExists, "user already has an account of this type ")
			}
		}
		return nil, status.Errorf(codes.Canceled, "account could not be created %s: ", err)
	}

	createAccountRes := convertAccount(account)

	return createAccountRes, nil

}

func ValidateCreateAccount(req *pb.CreateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, FieldViolation("currency", err))
	}

	if err := val.ValidateUsername(req.GetOwner()); err != nil {
		violations = append(violations, FieldViolation("owner", err))
	}

	return violations
}
