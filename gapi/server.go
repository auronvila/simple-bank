package gapi

import (
	"fmt"
	db "github.com/auronvila/simple-bank/db/sqlc"
	accountPb "github.com/auronvila/simple-bank/pb/account"
	userPb "github.com/auronvila/simple-bank/pb/user"
	"github.com/auronvila/simple-bank/token"
	"github.com/auronvila/simple-bank/util"
	"github.com/auronvila/simple-bank/worker"
)

type Server struct {
	userPb.UnimplementedUsersServer
	accountPb.UnimplementedAccountsServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new grpc server
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
