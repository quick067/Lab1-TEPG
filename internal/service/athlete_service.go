package service

import (
	"fmt"
	"training-system/internal/models"
)

type AthleteRepo interface {
	GetScheduleItems(userID uint, startDate, endDate string) ([]models.TrainingScheduleItem, error)
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