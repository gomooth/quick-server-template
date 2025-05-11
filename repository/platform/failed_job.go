package platform

import (
	"time"

	"gorm.io/gorm"
)

type FailedJob struct {
	gorm.Model

	JobName     string
	JobArgs     string
	Payload     string
	Errors      string
	Handled     bool
	HandledAt   *time.Time
	Compensated bool
}
