package routes

import "github.com/gin-gonic/gin"

func TransferRoutes(router *gin.Engine, createTransfer gin.HandlerFunc) {
	router.POST("/create-transfer", createTransfer)
}
