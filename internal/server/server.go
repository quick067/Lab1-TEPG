package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"training-system/internal/config"
)

type Server struct{
	engine *gin.Engine
	db *gorm.DB
	cfg config.Config
	server *http.Server
}

func NewServer(db *gorm.DB, cfg config.Config) *Server {
	engine := gin.Default()
	engine.Use(gin.Recovery())
	return &Server{
		engine: engine,
		db: db,
		cfg: cfg,
	}
}

func (s *Server) RunServer() error {
	newServ := http.Server{
		Addr: s.cfg.ServerPort,
		Handler: s.engine,
	}
	s.server = &newServ
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed{
		return err
	}
	return nil
}