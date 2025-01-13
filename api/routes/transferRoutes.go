package routes

import "github.com/gin-gonic/gin"

func TransferRoutes(router *gin.Engine, authRouter gin.IRoutes, createTransfer gin.HandlerFunc) {
	authRouter.POST("/create-transfer", createTransfer)
}
