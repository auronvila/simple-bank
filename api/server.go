package api

import (
	"github.com/auronvila/simple-bank/api/routes"
	simplebank "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  simplebank.Store
	router *gin.Engine
}

func NewServer(store simplebank.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	routes.AccountRoutes(
		router,
		server.CreateAccount,
		server.GetAccountById,
		server.ListAccounts,
	)

	routes.TransferRoutes(router, server.CreateTransfer)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
