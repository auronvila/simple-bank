package routes

import (
	"github.com/gin-gonic/gin"
)

func AccountRoutes(router *gin.Engine, authRouter gin.IRoutes, createAccount gin.HandlerFunc, getAccountById gin.HandlerFunc, listAccounts gin.HandlerFunc) {
	authRouter.POST("/accounts", createAccount)
	authRouter.GET("/account/:id", getAccountById)
	authRouter.GET("/accounts", listAccounts)
}
