package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// create HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	//调用的函数，参数必须有*gin.Context
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccount)
	server.router = router
	return server
}

// 开启监听
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// 很多地方都会用到
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
