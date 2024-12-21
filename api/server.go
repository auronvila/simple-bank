package api

import (
	"github.com/auronvila/simple-bank/api/routes"
	simplebank "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  simplebank.Store
	router *gin.Engine
}

func NewServer(store simplebank.Store) *Server {
	server := &Server{store: store}
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
	routes.UserRoutes(router, server.CreateUser, server.GetUserByUsername)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
