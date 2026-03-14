package storage

import (
	"fmt"
	"training-system/internal/models"

	"gorm.io/gorm"
)

type AthleteRepo struct {
	db *gorm.DB
}

func NewAthleteRepo(db *gorm.DB) *AthleteRepo {
	return &AthleteRepo{
		db: db,
	}
}

func (ar *AthleteRepo) GetScheduleItems(userID uint, startDate, endDate string) ([]models.TrainingScheduleItem, error) {
	var result []models.TrainingScheduleItem
	query := ar.db.Table("trainings").
		Select("title, scheduled_at, planned_duration").
		Where("coach_id = ?", 1)

	if startDate != "" && endDate != "" {
		query = query.Where("scheduled_at BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("error selecting schedule: %w", err)
	}

	return result, nil
}

func (ar *AthleteRepo) GetProgressItem(userID uint) ([]models.TrainingProgressItem, error) {
	var result []models.TrainingProgressItem

	query := ar.db.Table("trainings").
		Select(`
			trainings_logs.status, 
			trainings_logs.actual_duration, 
			trainings_logs.comment, 
			trainings.scheduled_at, 
			training_types.name
		`).
		Joins("INNER JOIN trainings_logs ON trainings.id = trainings_logs.training_id").
		Joins("INNER JOIN training_types ON trainings.task_type_id = training_types.id").
		Where("trainings_logs.athlete_id = ?", userID)

	query = query.Order("trainings.scheduled_at DESC")

	if err := query.Find(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch athlete progress: %w", err)
	}

	return result, nil
}

func (ar *AthleteRepo) PostHealthReport(userID uint, healthReport models.CreateHealthReportRequest) error {
	data := map[string]interface{}{
		"user_id": userID,
		"note": healthReport.Note,
	}

	err := ar.db.Table("health_reports").Create(data).Error
	if err != nil {
		return fmt.Errorf("error inserting health report: %w", err)
	}

	return nil
} 