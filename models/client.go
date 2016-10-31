package models

import (
	"github.com/cSploit/daemon/models/internal"
)

// A wifi client ( courtesy of aircrack )
type Client struct {
	internal.Base
	// MAC address
	Station string `json:"station"`
	Power   int    `json:"power"`
	Packets int    `json:"packets"`
	Bssid   string `json:"bssid"`
	Probed  string `json:"probed_essids"`

	Iface   Iface `json:"-"`
	IfaceId uint  `json:"-"`
	Jobs    []Job `json:"-" gorm:"many2many:job_clients"`
}

// DEAUTH infinitely the Client
func (c *Client) Deauth() (j Job, e error) {
	pj, e := CreateProcessJob("aireplay-ng", "-0", "0", "-a", c.Station, "-d", c.Bssid, c.Iface.Name)

	if e != nil {
		j = pj.Job
		internal.Db.Model(&j).Update("Name", "Deauth ["+c.Station+"]")
		internal.Db.Model(&j).Association("job_clients").Append(c)
		internal.Db.Model(&j).Association("job_ifaces").Append(&(c.Iface))
	}

	return
}
