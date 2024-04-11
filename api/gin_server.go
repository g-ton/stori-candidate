package api

import (
	db "github.com/g-ton/stori-candidate/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	Store  db.Store
	Router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store db.Store) *Server {
	return &Server{
		Store:  store,
		Router: gin.New(),
	}
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
