package api

import (
	simplebank "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *simplebank.Store
	router *gin.Engine
}

func NewServer(store *simplebank.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccountById)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
