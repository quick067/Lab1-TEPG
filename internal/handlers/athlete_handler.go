package handlers

import (
	"net/http"
	"training-system/internal/models"

	"github.com/gin-gonic/gin"
)

type AthleteService interface {
	GetSchedule(userID uint, startDate, endDate string) ([]models.TrainingScheduleItem, error)
	GetProgress(userID uint) ([]models.TrainingProgressItem, error)
	PostReport(userID uint, req models.CreateHealthReportRequest) error
}

type AthleteHandler struct {
	service AthleteService
}

func NewAthleteHandler(service AthleteService) *AthleteHandler {
	return &AthleteHandler{
		service: service,
	}
}

func (ah *AthleteHandler) GetScheduleHandler(c *gin.Context) {
	userID := uint(1) // hardcode (todo: add authentification)

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	schedule, err := ah.service.GetSchedule(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"schedule": schedule})
}

func (ah *AthleteHandler) GetProgressHandler(c *gin.Context) {
	userID := uint(1)

	progress, err := ah.service.GetProgress(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
    c.JSON(http.StatusOK, gin.H{"progress": progress})
}

func (ah *AthleteHandler) ReportHealthHandler(c *gin.Context) {
    userID := uint(1)

	var req models.CreateHealthReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return 
	}
	if err := ah.service.PostReport(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}
	c.JSON(http.StatusCreated, gin.H{"message": "health report created successfully"})
}
