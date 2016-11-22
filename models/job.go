package models

import (
	"database/sql/driver"
	"github.com/cSploit/daemon/models/internal"
	"time"
)

func init() {
	internal.RegisterModels(&Job{})
	internal.RegisterJoinTables("job_hosts", "job_aps", "job_networks", "job_clients", "job_ifaces")
}

type (
	JobKind int64

	// A running task
	Job struct {
		internal.Base
		FinishedAt *time.Time `json:"finished_at"`
		Name       string     `json:"name"`
		Type       JobKind    `json:"type"`

		// affected entities
		Aps      []AP      `json:"-" gorm:"many2many:job_aps"`
		Clients  []Client  `json:"-" gorm:"many2many:job_clients"`
		Hosts    []Host    `json:"-" gorm:"many2many:job_hosts"`
		Networks []Network `json:"-" gorm:"many2many:job_networks"`
		Ifaces   []Iface   `json:"-" gorm:"many2many:job_ifaces"`

		// concrete jobs
		Radar   *RadarJob   `json:"-"`
		Process *ProcessJob `json:"-"`
		//TODO: DiscoveryJob, MonitorJob
	}
)

var jobKindNames = map[JobKind]string{}

func registerJobKind(kind JobKind, name string /*, ViewHandlerFunction*/) {
	if _, ok := jobKindNames[kind]; ok {
		panic("job kind already registered: " + name)
	}
	jobKindNames[kind] = name
}

func (k JobKind) String() string {
	return jobKindNames[k]
}

// used for json serialization
func (k JobKind) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

// DB deserialization
func (k *JobKind) Scan(value interface{}) error {
	*k = JobKind(value.(int64))
	return nil
}

// DB serialization
func (k JobKind) Value() (driver.Value, error) {
	return int64(k), nil
}

func (j *Job) Is(kind JobKind) bool {
	return j.Type == kind
}

func FindJob(id uint) (j *Job, e error) {
	j = &Job{}
	e = internal.Db.Find(j, id).Error
	return
}
