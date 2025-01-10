package api

import (
	"fmt"
	"github.com/auronvila/simple-bank/api/routes"
	db "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/auronvila/simple-bank/token"
	"github.com/auronvila/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	routes.AccountRoutes(
		router,
		server.CreateAccount,
		server.GetAccountById,
		server.ListAccounts,
	)
	routes.TransferRoutes(router, server.CreateTransfer)
	routes.UserRoutes(router, server.CreateUser, server.GetUserByUsername, server.loginUser)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
