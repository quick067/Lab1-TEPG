package models

import (
	"time"
)

type Status string

const (
	Done   Status = "Done"
	Missed Status = "Missed"
)

type TrainingScheduleItem struct {
	Title           string    `json:"title"`
	ScheduledAt     time.Time `json:"scheduled_at"`
	PlannedDuration uint      `json:"planned_duration"`
	Status          Status    `json:"status"`
	ActualDuration  uint      `json:"actual_duration"`
	Comment         string    `json:"comment"`
}

type CreateHealthReportRequest struct {
	Note string `json:"note"`
}

type TrainingProgressItem struct {
	Status Status `json:"status"`
	ActualDuration uint `json:"actual_duration"`
	Comment string `json:"comment"`
	ScheduledAt time.Time `json:"scheduled_at"`
	TrainingTypeName string `json:"name"`
}