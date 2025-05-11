package platform

import "gorm.io/gorm"

type FailedListener struct {
	gorm.Model

	ConsumeGroup  string
	Topic         string
	Msg           string
	FailedPayload string
	FailedReason  string
}
