package storage

import (
	"fmt"
	"training-system/internal/models"

	"time"

	"gorm.io/gorm"
)

type CoachRepo struct {
	db *gorm.DB
}

func NewCoachRepo(db *gorm.DB) *CoachRepo {
	return &CoachRepo{
		db: db,
	}
}

type dbTraining struct {
	CoachID         uint      `gorm:"column:coach_id"`
	TaskTypeID      uint      `gorm:"column:task_type_id"`
	Title           string    `gorm:"column:title"`
	Description     string    `gorm:"column:description"`
	PlannedDuration uint      `gorm:"column:planned_duration"`
	ScheduledAt     time.Time `gorm:"column:scheduled_at"`
}

func (cr *CoachRepo) CreateTraining(userID uint, training models.CreateTrainingRequest) error {
	row := dbTraining{
		CoachID:         userID,
		TaskTypeID:      training.TrainingTypeId,
		Title:           training.Title,
		Description:     training.Description,
		PlannedDuration: training.PlannedDuration,
		ScheduledAt:     training.ScheduledAt,
	}

	err := cr.db.Table("trainings").Create(&row).Error
	if err != nil {
		return fmt.Errorf("error inserting training data: %w", err)
	}

	return nil
}

func (cr *CoachRepo) GetTrainingsAnalytics(userID uint, startDate, endDate string) ([]models.TrainingsAnalyticsItem, error) {
	result := []models.TrainingsAnalyticsItem{}
	query := cr.db.Table("trainings").
        Select("trainings.title, trainings.scheduled_at, trainings.planned_duration, trainings_logs.actual_duration, COALESCE(trainings_logs.comment, trainings.description) AS comment").
        Joins("LEFT JOIN trainings_logs ON trainings.id = trainings_logs.training_id").
        Where("trainings.coach_id = ?", userID)

	if len(startDate) != 0 && len(endDate) != 0 {
		query = query.Where("scheduled_at BETWEEN ? AND ?", startDate, endDate)
	}

	query = query.Scan(&result)

	if query.Error != nil {
		return nil, fmt.Errorf("error selecting values: %w", query.Error)
	}

	return result, nil
}

func (cr *CoachRepo) AddMember(athleteID uint) error {
	query := cr.db.Table("users").Where("id = ?", athleteID).Update("is_active", true)

	if query.Error != nil {
		return fmt.Errorf("error updating values: %w", query.Error)
	}

	if query.RowsAffected == 0 {
		return fmt.Errorf("athlete with id %d not found", athleteID)
	}
	return nil
}

func (cr *CoachRepo) DeleteMember(athleteID uint) error {
	query := cr.db.Table("users").Where("id = ?", athleteID).Update("is_active", false)

	if query.Error != nil {
		return fmt.Errorf("error updating values: %w", query.Error)
	}

	if query.RowsAffected == 0 {
		return fmt.Errorf("athlete with id %d not found", athleteID)
	}
	return nil
}

func (cr *CoachRepo) UpdateTraining(trainingID uint, req models.UpdateTrainingLogRequest) error {
	data := map[string]interface{}{
		"status":          req.Status,
		"actual_duration": req.ActualDuration,
		"comment":         req.Comment,
	}

	query := cr.db.Table("trainings_logs").Where("id = ?", trainingID).Updates(data)

	if query.Error != nil {
		return fmt.Errorf("error updating values: %w", query.Error)
	}

	if query.RowsAffected == 0 {
		return fmt.Errorf("training log with id %d not found", trainingID)
	}
	return nil
}

func (cr *CoachRepo) GetTeamMembers() ([]models.User, error) {
	var athletes []models.User
	
	err := cr.db.Table("users").Where("role = ?", "athlete").Find(&athletes).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching team members: %w", err)
	}
	
	return athletes, nil
}