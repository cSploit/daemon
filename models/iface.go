package models

import "github.com/cSploit/daemon/models/internal"

func init() {
	internal.RegisterModels(&Iface{})
}

// A network interface
type Iface struct {
	internal.Base
	Name string `json:"name"`

	Aps     []AP     `json:"-"`
	Clients []Client `json:"-"`
	Jobs    []Job    `json:"-" gorm:"many2many:job_ifaces"`
}
