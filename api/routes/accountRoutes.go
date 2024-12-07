package routes

import (
	"github.com/gin-gonic/gin"
)

func AccountRoutes(router *gin.Engine, createAccount gin.HandlerFunc, getAccountById gin.HandlerFunc, listAccounts gin.HandlerFunc) {
	router.POST("/accounts", createAccount)
	router.GET("/account/:id", getAccountById)
	router.GET("/accounts", listAccounts)
}
