package models

import (
	"github.com/cSploit/daemon/models/internal"
	"time"
)

// A wifi client ( courtesy of aircrack )
type Client struct {
	internal.Base
	// MAC address
	First   time.Time `json:"first_seen"`
	Last    time.Time `json:"last_seen"`
	Station string    `json:"station"`
	Power   int       `json:"power"`
	Packets int       `json:"packets"`
	Bssid   string    `json:"bssid"`
	Probed  string    `json:"probed_essids"`

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

func FindClient(id uint) (c *Client, e error) {
	c = &Client{}
	e = internal.Db.Find(c, id).Error
	return
}

func FindClientByMac(mac_addr string) (c *Client, e error) {
	c = &Client{}
	e = internal.Db.Where("station = ?", mac_addr).Find(c).Error
	return
}
