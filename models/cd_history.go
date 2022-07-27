package models

import (
	"gorm.io/gorm"
	"time"
)

type CDHistory struct {
	gorm.Model
	App string `gorm:"index:app_env;not null"`
	Env string `gorm:"index:app_env;not null"`
	CommitID string `gorm:"index;not null"`
	Proj string `gorm:"not null"`
	GitUrl string `gorm:"not null"`
	Branch string `gorm:"not null"`
	ImageTag string `gorm:"not null"`
	ImageUrl string `gorm:"not null"`
	DeployTime time.Time `gorm:"not null"`
}

type CDStatus struct {
	gorm.Model
	App string `gorm:"index:app_env;not null"`
	Env string `gorm:"index:app_env;not null"`
	CommitID string `gorm:"index;not null"`
	Proj string `gorm:"not null"`
	GitUrl string `gorm:"not null"`
	Branch string `gorm:"not null"`
	ImageTag string `gorm:"not null"`
	ImageUrl string `gorm:"not null"`
}
