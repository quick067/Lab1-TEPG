package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"training-system/internal/config"
	"training-system/internal/handlers"
	"training-system/internal/service"
	"training-system/internal/storage"
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
	AtheleteRepo := storage.NewAthleteRepo(s.db)
	AtheleteService := service.NewAthleteService(AtheleteRepo)
	AtheleteHandlers := handlers.NewAthleteHandler(AtheleteService)

	v1 := s.engine.Group("v1")

	athelete := v1.Group("athelete")
	{
		athelete.GET("/schedule", AtheleteHandlers.GetScheduleHandler) 
		athelete.GET("/progress", AtheleteHandlers.GetProgressHandler)
		athelete.POST("/health-report", AtheleteHandlers.ReportHealtsHandler)
	}

	coach := v1.Group("coach")
	{
		coach.POST("/training", ) // CreateTrainingHandler
		coach.POST("/team/members", ) //AddMemberHandler
		coach.DELETE("/team/members/:athlete_id", ) //DeleteMemberHandler
		coach.PUT("/training/logs/:id", ) //UpdateTrainingLogs
		coach.GET("/analytics", ) //GetAnalytics
	}

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