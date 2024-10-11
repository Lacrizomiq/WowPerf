package models

import (
	"time"

	"gorm.io/gorm"
)

type UpdateState struct {
	gorm.Model
	LastUpdateTime time.Time
}
