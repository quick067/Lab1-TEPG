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
	result := []models.TrainingScheduleItem{}
	query := ar.db.Table("trainings").
	Select("trainings.title, trainings.scheduled_at, trainings.planned_duration, trainings_logs.status, trainings_logs.actual_duration, trainings_logs.training_comment").
	Joins("JOIN trainings_logs ON trainings.id=trainings_logs.training_id").
	Where("trainings_logs.athlete_id = ?", userID)

	if len(startDate) != 0 && len(endDate) != 0 {
		query = query.Where("trainings.scheduled_at BETWEEN ? AND ?", startDate, endDate)
	}

	query = query.Scan(&result)

	if query.Error != nil {
		return nil, fmt.Errorf("error selecting values: %w", query.Error)
	}

	return result, nil
}

func (ar *AthleteRepo) GetProgressItem(userID uint) ([]models.TrainingProgressItem, error){
	result := []models.TrainingProgressItem{}
	query := ar.db.Table("trainings").
	Select("status, actual_duration, scheduled_at, name").
	Joins("JOIN trainings_logs ON trainings.id = trainings_logs.training_id").
	Joins("JOIN training_types ON trainings.task_type_id = training_types.id").
	Where("trainings_logs.athlete_id = ?", userID)

	query = query.Scan(&result)

	if query.Error != nil {
		return nil, fmt.Errorf("error selecting values: %w", query.Error)
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