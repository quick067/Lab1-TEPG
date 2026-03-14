package models

import "time"

type CreateTrainingRequest struct {
	Title           string    `json:"title" binding:"required,min=3"`
	TrainingTypeId  uint      `json:"training_type_id" binding:"required,gt=0"`
	ScheduledAt     time.Time `json:"scheduled_at" binding:"required"`
	PlannedDuration uint      `json:"planned_duration" binding:"required,gt=0"`
	Description     string    `json:"description" binding:"omitempty,max=1000"`
}

type AddMemberRequest struct {
	MemberID uint `json:"member_id" binding:"required,gt=0"`
}

type TrainingsAnalyticsItem struct {
	ID              uint      `json:"id"`
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

type UpdateLogTrainingRequest struct {
	ActualDuration uint   `json:"actual_duration" binding:"omitempty,gt=0"`
	Comment        string `json:"comment" binding:"omitempty,max=1000"`
	Status         string `json:"status" binding:"required,oneof=Done Missed"`
}
