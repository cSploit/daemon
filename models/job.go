package models

import (
	"time"
)

type (
	Job struct {
		base
		FinishedAt *time.Time `json:"finished_at"`
		Name       string     `json:"name"`
	}
)

func FindJob(id uint) (j *Job, e error) {
	j = &Job{}
	e = db.Find(j, id).Error
	return
}
