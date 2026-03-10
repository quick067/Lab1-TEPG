package service

import (
	"fmt"
	"training-system/internal/models"
)

type AthleteRepo interface {
	GetScheduleItems(userID uint, startDate, endDate string) ([]models.TrainingScheduleItem, error)
	GetProgressItem(userID uint) ([]models.TrainingProgressItem, error)
	PostHealthReport(userID uint, healthReport models.CreateHealthReportRequest) error
}

type AthleteService struct {
	repo AthleteRepo
}

func NewAthleteService(repo AthleteRepo) *AthleteService {
	return &AthleteService{
		repo: repo,
	}
}

func (as *AthleteService) GetSchedule(userID uint, startDate, endDate string) ([]models.TrainingScheduleItem, error) {
	res, err := as.repo.GetScheduleItems(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule from repo: %w", err)
	}

	return res, nil
}

func (as *AthleteService) GetProgress(userID uint) ([]models.TrainingProgressItem, error) {
	res, err := as.repo.GetProgressItem(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress from repo: %w", err)
	}
	return res, nil
}

func (as *AthleteService) PostReport (userID uint, healthReport models.CreateHealthReportRequest) error {
	if err := as.repo.PostHealthReport(userID, healthReport); err != nil {
		return fmt.Errorf("failed to create health report in db: %w", err)
	}
	return nil
}