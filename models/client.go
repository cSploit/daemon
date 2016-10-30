package models

import "github.com/cSploit/daemon/models/internal"

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
