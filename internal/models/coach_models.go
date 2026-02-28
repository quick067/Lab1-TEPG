package models

import "time"

type CreateTrainingRequest struct {
	TrainingTypeId  uint      `json:"name"`
	ScheduledAt     time.Time `json:"scheduled_at"`
	PlannedDuration uint      `json:"planned_duration"`
	Description     string    `json:"description"`
}

type AddMemberRequest struct {
	MemberID uint `json:"id"`
}

type UpdateTrainingLogRequest struct {
	Status Status `json:"status"`
	ActualDuration uint `json:"actual_duration"`
	Comment string `json:"comment"`
}

type TrainingsAnalyticsItem struct {
	Title           string    `json:"title"`
	ScheduledAt     time.Time `json:"scheduled_at"`
	PlannedDuration uint      `json:"planned_duration"`
	ActualDuration  uint      `json:"actual_duration"`
	Comment         string    `json:"comment"`
}