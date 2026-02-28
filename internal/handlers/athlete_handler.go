package handlers

import (
	"net/http"
	"training-system/internal/models"

	"github.com/gin-gonic/gin"
)

type AthleteService interface {
	GetSchedule(userID uint, startDate, endDate string) ([]models.TrainingScheduleItem, error)
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
    c.JSON(http.StatusOK, gin.H{"todo": "todo"})
}

func (ah *AthleteHandler) ReportHealtsHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"todo": "todo"})
}
