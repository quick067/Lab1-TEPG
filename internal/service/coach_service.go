package service

import (
	"fmt"
	"training-system/internal/models"
)

type CoachRepo interface{
	GetTrainingsAnalytics(userID uint, startDate, endDate string) ([]models.TrainingsAnalyticsItem, error)
	CreateTraining(userID uint, training models.CreateTrainingRequest) error
	AddMember(athleteID uint) error
	DeleteMember(athleteID uint) error 
	UpdateTraining(trainingID uint, req models.UpdateTrainingLogRequest) error
	GetTeamMembers() ([]models.User, error)
}

type CoachService struct {
	repo CoachRepo
}

func NewCoachService(repo CoachRepo) *CoachService {
	return &CoachService{
		repo: repo,
	}
}

func (cs *CoachService) GetTrainingsAnalytics (userID uint, startDate, endDate string) ([]models.TrainingsAnalyticsItem, error) {
	res, err := cs.repo.GetTrainingsAnalytics(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics from repo: %w", err)
	}

	return res, err
}

func (cs *CoachService) CreateTraining (userID uint, training models.CreateTrainingRequest) error {
	err := cs.repo.CreateTraining(userID, training)
	if err != nil {
		return fmt.Errorf("error creating training: %w", err)
	}
	return nil
}

func (cs *CoachService) UpdateTraining(trainingID uint, req models.UpdateTrainingLogRequest) error {
	err := cs.repo.UpdateTraining(trainingID, req)
	if err != nil {
		return fmt.Errorf("error updating values: %w", err)
	}
	return nil
}
func (cs *CoachService) DeleteTeamMember(athleteID uint) error {
	err := cs.repo.DeleteMember(athleteID)
	if err != nil {
		return fmt.Errorf("error deleting member: %w", err)
	}
	return nil 
}

func (cs *CoachService)AddTeamMember(MemberID uint) error {
	err := cs.repo.AddMember(MemberID)
	if err != nil {
		return fmt.Errorf("error inserting member: %w", err)
	}
	return nil
}

func (cs *CoachService) GetTeamMembers() ([]models.User, error) {
	athletes, err := cs.repo.GetTeamMembers()
	if err != nil {
		return nil, fmt.Errorf("failed to get team members from repo: %w", err)
	}
	return athletes, nil
}