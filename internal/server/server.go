package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func New() *Server {
	router := gin.Default()

	srv := &Server{
		router: router,
	}

	srv.setupRoutes()

	return srv
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) setupRoutes() {
	setupRoutes(s.router)
}
