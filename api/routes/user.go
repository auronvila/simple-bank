package routes

import "github.com/gin-gonic/gin"

func UserRoutes(router *gin.Engine, createUser gin.HandlerFunc, getUserByUsername gin.HandlerFunc, loginUser gin.HandlerFunc) {
	router.POST("/create-user", createUser)
	router.GET("/user/:username", getUserByUsername)
	router.POST("/user/login-user", loginUser)
}
