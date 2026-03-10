package models

import "time"

type CreateTrainingRequest struct {
	Title           string    `json:"title"`
	TrainingTypeId  uint      `json:"training_type_id"`
	ScheduledAt     time.Time `json:"scheduled_at"`
	PlannedDuration uint      `json:"planned_duration"`
	Description     string    `json:"description"`
}

type AddMemberRequest struct {
	MemberID uint `json:"member_id"`
}

type UpdateTrainingLogRequest struct {
	Status         Status `json:"status"`
	ActualDuration uint   `json:"actual_duration"`
	Comment        string `json:"comment"`
}

type TrainingsAnalyticsItem struct {
	Title           string    `json:"title"`
	ScheduledAt     time.Time `json:"scheduled_at"`
	PlannedDuration uint      `json:"planned_duration"`
	ActualDuration  uint      `json:"actual_duration"`
	Comment         string    `json:"comment"`
}

type User struct {
	ID        uint   `gorm:"primaryKey;column:id"`
	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Email     string `gorm:"column:email"`
	Role      string `gorm:"column:role"`
	IsActive  bool   `gorm:"column:is_active"`
}