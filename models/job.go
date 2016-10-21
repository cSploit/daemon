package models

import (
	"time"
)

type (
	Job struct {
		ID         uint       `gorm:"primary_key" json:"id"`
		CreatedAt  time.Time  `json:"created_at"`
		UpdatedAt  time.Time  `json:"updated_at"`
		FinishedAt *time.Time `json:"finished_at"`
		Name       string     `json:"name"`
	}
)

func FindJob(id uint) (j *Job, e error) {
	j = &Job{}
	e = db.Find(j, id).Error
	return
}
