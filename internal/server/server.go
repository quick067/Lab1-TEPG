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
	engine.LoadHTMLGlob("templates/*")
	engine.Use(gin.Recovery())
	return &Server{
		engine: engine,
		db: db,
		cfg: cfg,
	}
}

func (s *Server) RunServer() error {
	AthleteRepo := storage.NewAthleteRepo(s.db)
	AtheleteService := service.NewAthleteService(AthleteRepo)
	AthleteHandlers := handlers.NewAthleteHandler(AtheleteService)

	CoachRepo := storage.NewCoachRepo(s.db)
	CoachService := service.NewCoachService(CoachRepo)
	CoachHandler := handlers.NewCoachHandler(CoachService)

	api := s.engine.Group("/api/v1")
	web := s.engine.Group("/web")
	{
		web.GET("/coach/dashboard", CoachHandler.GetDashboardView)
		web.GET("/coach/team", CoachHandler.GetTeamView)
		web.GET("/athlete/dashboard", AthleteHandlers.GetDashboardView)
   		web.GET("/athlete/health", AthleteHandlers.GetHealthReportView)
    	web.POST("/athlete/health", AthleteHandlers.CreateHealthReportWeb)
	}


	athlete := api.Group("athlete")
	{
		athlete.GET("/schedule", AthleteHandlers.GetScheduleHandler) 
		athlete.GET("/progress", AthleteHandlers.GetProgressHandler)
		athlete.POST("/health-report", AthleteHandlers.ReportHealthHandler)
	}

	coach := api.Group("coach")
	{
		coach.POST("/training", CoachHandler.CreateTrainingHandler) 
		coach.POST("/team/members", CoachHandler.AddMemberHandler)
		coach.DELETE("/team/members/:athlete_id", CoachHandler.DeleteMemberHandler)
		coach.PUT("/training/logs/:id", CoachHandler.UpdateTrainingLogs)
		coach.GET("/analytics", CoachHandler.GetAnalyticsHandler) 
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