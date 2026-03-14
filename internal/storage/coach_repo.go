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
		Select(`
			trainings.id, 
			trainings.title, 
			trainings.scheduled_at, 
			trainings.planned_duration, 
			trainings_logs.actual_duration, 
			COALESCE(trainings_logs.comment, trainings.description) AS comment
		`).
		Joins("LEFT JOIN trainings_logs ON trainings.id = trainings_logs.training_id").
		Where("trainings.coach_id = ?", userID)

	if len(startDate) != 0 && len(endDate) != 0 {
		query = query.Where("scheduled_at BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Order("trainings.scheduled_at DESC").Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("error selecting values: %w", err)
	}

	return result, nil
}

func (cr *CoachRepo) AddMember(athleteID uint) error {
	return cr.updateUserStatus(athleteID, true)
}

func (cr *CoachRepo) DeleteMember(athleteID uint) error {
	return cr.updateUserStatus(athleteID, false)
}

func (cr *CoachRepo) updateUserStatus(athleteID uint, active bool) error {
	res := cr.db.Table("users").Where("id = ?", athleteID).Update("is_active", active)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("user %d not found", athleteID)
	}
	return nil
}

func (cr *CoachRepo) UpdateTraining(trainingID uint, req models.UpdateLogTrainingRequest) error {
	athleteID := uint(2)

	var count int64
	cr.db.Table("trainings_logs").
		Where("training_id = ? AND athlete_id = ?", trainingID, athleteID).
		Count(&count)

	if count > 0 {
		data := map[string]interface{}{
			"status":          req.Status,
			"actual_duration": req.ActualDuration,
			"comment":         req.Comment,
		}

		err := cr.db.Table("trainings_logs").
			Where("training_id = ? AND athlete_id = ?", trainingID, athleteID).
			Updates(data).Error

		if err != nil {
			return fmt.Errorf("error updating training log: %w", err)
		}
		return nil
	}

	newData := map[string]interface{}{
		"training_id":     trainingID,
		"athlete_id":      athleteID,
		"status":          req.Status,
		"actual_duration": req.ActualDuration,
		"comment":         req.Comment,
	}

	err := cr.db.Table("trainings_logs").Create(newData).Error
	if err != nil {
		return fmt.Errorf("error creating training log: %w", err)
	}

	return nil
}

func (cr *CoachRepo) GetTeamMembers() ([]models.User, error) {
	var athletes []models.User

	err := cr.db.Table("users").
		Where("role = ?", "athlete").
		Find(&athletes).
		Error

	if err != nil {
		return nil, fmt.Errorf("error fetching team members: %w", err)
	}

	return athletes, nil
}

func (cr *CoachRepo) GetAllTeamTrainings(coachID uint) ([]models.TrainingScheduleItem, error) {
	var result []models.TrainingScheduleItem
	err := cr.db.Table("trainings").
		Select("title, scheduled_at, planned_duration").
		Where("coach_id = ?", coachID).
		Order("scheduled_at DESC").
		Find(&result).Error
	return result, err
}
