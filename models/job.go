package models

import (
	"github.com/cSploit/daemon/models/internal"
	"time"
)

func init() {
	internal.RegisterModels(&Job{})
	internal.RegisterJoinTables("job_hosts", "job_aps", "job_networks", "job_clients")
}

type (
	Job struct {
		internal.Base
		FinishedAt *time.Time `json:"finished_at"`
		Name       string     `json:"name"`
		Aps        []AP       `json:"aps" gorm:"many2many:job_aps"`
		Clients    []Client   `json:"clients" gorm:"many2many:job_clients"`
		Hosts      []Host     `json:"hosts" gorm:"many2many:job_hosts"`
		Networks   []Network  `json:"networks" gorm:"many2many:job_networks"`
	}
)

func FindJob(id uint) (j *Job, e error) {
	j = &Job{}
	e = internal.Db.Find(j, id).Error
	return
}
