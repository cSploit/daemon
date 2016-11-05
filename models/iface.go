package models

import (
	"github.com/cSploit/daemon/models/internal"
	"io/ioutil"
	"os"
)

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

func (iface *Iface) StartDiscovery() (d *DiscoveryJob, e error) {
	dir, e := ioutil.TempDir("", "airodump-")

	if e != nil {
		return
	}

	pj, e := CreateProcessJob("airodump-ng", "--write", os.TempDir()+"/discovery", "--output-format", "csv", "--wps", iface.Name)

	if e != nil {
		os.Remove(dir)
		return
	}

	d = &DiscoveryJob{}
	d.Dir = dir
	d.Job = pj.Job

	e = internal.Db.Save(d).Error
	return
}
