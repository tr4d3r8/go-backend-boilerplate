package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/tr4d3r8/go-backend-boilerplate/db/sqlc"
)

//  serve all http requests
type Server struct {
	store  db.Store
	router *gin.Engine
}

// create  new server instance
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/account", server.listAccount)
	// add routes to router

	server.router = router
	return server
}

// runs the http server on specific address and starts listening for requests
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
