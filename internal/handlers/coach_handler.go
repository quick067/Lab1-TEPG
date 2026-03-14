package handlers

import (
	"net/http"
	"strconv"
	"training-system/internal/models"

	"github.com/gin-gonic/gin"
)

type CoachService interface {
	GetTrainingsAnalytics(userID uint, startDate, endDate string) ([]models.TrainingsAnalyticsItem, error)
	CreateTraining(userID uint, training models.CreateTrainingRequest) error
	UpdateTraining(trainingID uint, req models.UpdateLogTrainingRequest) error
	DeleteTeamMember(athleteID uint) error
	AddTeamMember(MemberID uint) error
	GetTeamMembers() ([]models.User, error)
	GetAllTeamTrainings(coachID uint) ([]models.TrainingScheduleItem, error)
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

func (ch *CoachHandler) GetAnalyticsHandler(c *gin.Context) {
	coachID := uint(1) // hardcode (todo: add authentification)

	analytics, err := ch.service.GetTrainingsAnalytics(coachID, "", "")
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
	trainingIDStr := c.Param("id")
	trainingID, err := strconv.Atoi(trainingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невірний ID тренування"})
		return
	}

	var req models.UpdateLogTrainingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ch.service.UpdateTraining(uint(trainingID), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Результат збережено"})
}

func (ch *CoachHandler) GetDashboardView(c *gin.Context) {
	c.HTML(http.StatusOK, "coach_dashboard.html", gin.H{
		"Title": "Кабінет тренера",
	})
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