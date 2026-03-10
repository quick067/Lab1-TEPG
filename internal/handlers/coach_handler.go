package handlers

import (
	"net/http"
	"strconv"
	"time"
	"training-system/internal/models"

	"github.com/gin-gonic/gin"
)

type CoachService interface {
	GetTrainingsAnalytics(userID uint, startDate, endDate string) ([]models.TrainingsAnalyticsItem, error)
	CreateTraining(userID uint, training models.CreateTrainingRequest) error
	UpdateTraining(trainingID uint, req models.UpdateTrainingLogRequest) error
	DeleteTeamMember(athleteID uint) error
	AddTeamMember(MemberID uint) error
	GetTeamMembers() ([]models.User, error)
}

type CoachHandler struct {
	service CoachService
}

func NewCoachHandler(service CoachService) *CoachHandler {
	return &CoachHandler{
		service: service,
	}
}

func (ch *CoachHandler) CreateTrainingHandler(c *gin.Context) {
	userID := uint(1)

	var req models.CreateTrainingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}
	if err := ch.service.CreateTraining(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "training created successfully"})
}

func (ch *CoachHandler) GetAnalyticsHandler (c *gin.Context) {
	userID := uint(1) // hardcode (todo: add authentification)

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	analytics, err := ch.service.GetTrainingsAnalytics(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"analytics": analytics})
}

func (ch *CoachHandler) AddMemberHandler(c *gin.Context) {
	var req models.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}
	if err := ch.service.AddTeamMember(req.MemberID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "member added successfully"})
}

func (ch *CoachHandler) DeleteMemberHandler(c *gin.Context) {	
	idInt, err := strconv.Atoi(c.Param("athlete_id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid athlete ID"})
		return
	}
	athleteID := uint(idInt)

	if err := ch.service.DeleteTeamMember(athleteID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member deleted successfully"})
}

func (ch *CoachHandler) UpdateTrainingLogs(c *gin.Context) {
	idInt, err := strconv.Atoi(c.Param("training_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}
	trainingID := uint(idInt)
	var req models.UpdateTrainingLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := ch.service.UpdateTraining(trainingID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "training updated successfully"})
}


func (ch *CoachHandler) GetDashboardView(c *gin.Context) {
	userID := uint(1)

	analytics, err := ch.service.GetTrainingsAnalytics(userID, "", "")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "dashboard.html", gin.H{"Title": "Панель тренера", "Error": "Помилка завантаження даних"})
		return 
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Title": "Панель тренера",
		"Analytics": analytics,
	})
}

func (ch *CoachHandler) CreateTrainingWebHandler(c *gin.Context){
	userID := uint(1)

	title := c.PostForm("title")
	description := c.PostForm("description")

	trainingTypeID, _ := strconv.Atoi(c.PostForm("training_type_id"))
	PlannedDuration, _ := strconv.Atoi(c.PostForm("planned_duration"))

	scheduledAtStr := c.PostForm("scheduled_at")
	scheduledAt, _ := time.Parse("2006-01-02T15:04", scheduledAtStr)

	req := models.CreateTrainingRequest {
		Title: title,
		Description: description,
		TrainingTypeId: uint(trainingTypeID),
		PlannedDuration: uint(PlannedDuration),
		ScheduledAt: scheduledAt,
	}

	if err := ch.service.CreateTraining(userID, req); err != nil {
		c.HTML(http.StatusInternalServerError, "dashboard.html", gin.H{"title": "Панель тренера", "Error": err.Error()})
		return 
	}
	c.Redirect(http.StatusFound, "/web/coach/dashboard")
}

func (ch *CoachHandler) GetTeamView(c *gin.Context) {
	athletes, err := ch.service.GetTeamMembers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "team.html", gin.H{"Title": "Моя команда", "Error": "Помилка завантаження списку атлетів"})
		return
	}

	c.HTML(http.StatusOK, "team.html", gin.H{
		"Title":    "Керування командою",
		"Athletes": athletes,
	})
}