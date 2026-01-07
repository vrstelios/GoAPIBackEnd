package internal

import (
	"GoAPIBackEnd/config"
	"GoAPIBackEnd/internal/helpers"
	"GoAPIBackEnd/internal/http"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func New() *Server {
	helpers.SetJWTKey(config.GetConfig().JWT.Secret)
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
	http.SetupRoutes(s.router)
}
