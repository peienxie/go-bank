package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/peienxie/go-bank/db/sqlc"
)

// Server serves HTTP requests
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup its routing
func NewServer(store db.Store) *Server {
	server := &Server{
		store:  store,
		router: gin.Default(),
	}

	// initilizes routing
	server.initAccountRoutes()
	server.initTransferRoutes()

	return server
}

// Serve runs the http server on the provided address
func (s *Server) Serve(addr string) error {
	return s.router.Run(addr)
}
