package handlers

import (
	"net/http"
	"time"
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
	userID := uint(2) // hardcode

	now := time.Now()

	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startDate := startOfToday.Format("2006-01-02 15:04:05")

	endOfWindow := time.Date(now.Year(), now.Month(), now.Day()+7, 23, 59, 59, 0, now.Location())
	endDate := endOfWindow.Format("2006-01-02 15:04:05")

	schedule, err := ah.service.GetSchedule(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"schedule": schedule})
}

func (ah *AthleteHandler) GetProgressHandler(c *gin.Context) {
	userID := uint(2)

	progress, err := ah.service.GetProgress(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
    c.JSON(http.StatusOK, gin.H{"progress": progress})
}

func (ah *AthleteHandler) ReportHealthHandler(c *gin.Context) {
    userID := uint(2)

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


func (ah *AthleteHandler) GetDashboardView(c *gin.Context) {
    c.HTML(http.StatusOK, "athlete_dashboard.html", gin.H{
        "Title": "Кабінет атлета",
    })
}

func (ah *AthleteHandler) GetHealthReportView(c *gin.Context) {
    c.HTML(http.StatusOK, "health_report.html", gin.H{
        "Title": "Звіт про стан здоров'я",
    })
}

func (ah *AthleteHandler) CreateHealthReportWeb(c *gin.Context) {
	athleteID := uint(2)
	
	note := c.PostForm("note")

	req := models.CreateHealthReportRequest{
		Note: note,
	}

	if err := ah.service.PostReport(athleteID, req); err != nil {
		c.HTML(http.StatusInternalServerError, "health_report.html", gin.H{
			"Title": "Самопочуття", 
			"Error": "Помилка збереження звіту: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/web/athlete/dashboard")
}