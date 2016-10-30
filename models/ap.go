package models

import (
	"github.com/cSploit/daemon/models/internal"
	"github.com/cSploit/daemon/tools/aircrack/attacks"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func init() {
	internal.RegisterModels(&AP{})
}

// Access Point ( courtesy of aircrack )
type AP struct {
	internal.Base
	Bssid   string `json:"bssid"`
	Channel int    `json:"channel"`
	Speed   int    `json:"speed"`
	Privacy string `json:"privacy"`
	Cipher  string `json:"cipher"`
	Auth    string `json:"auth"`
	Power   int    `json:"power"`
	Beacons int    `json:"beacons"`
	IVs     int    `json:"ivs"`
	Lan     string `json:"lan_ip"`
	IdLen   int    `json:"id_len"`
	Essid   string `json:"essid"`
	Key     string `json:"key"`
	//Wps     bool   `json:"wps"`

	Iface   Iface `json:"-"`
	IfaceId uint  `json:"-"`
	Jobs    []Job `json:"-" gorm:"many2many:job_aps;"`
}

// DEAUTH infinitely the AP using broadcast address
func (a *AP) Deauth(iface string) (j Job, e error) {
	pj, e := CreateProcessJob("aireplay-ng", "-0", "0", "-a", a.Bssid, iface)

	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Name", "Deauth ["+a.Bssid+"]")
		db.Model(&j).Association("Aps").Append(a)
	}

	return
}

// Try a fake auth on the ap
func (a *AP) FakeAuth(iface string) (j Job, e error) {
	pj, e := CreateProcessJob("aireplay-ng", "-1", "0", "-a", a.Bssid, "-T", "1", iface)

	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Name", "FakeAuth ["+a.Bssid+"]")
		db.Model(&j).Association("Aps").Append(a)
	}

	return
}

// ARP replay!!
func (a *AP) ArpReplay(iface string) (attacks.Attack, error) {
	cmd := exec.Command("aireplay-ng", "-3", "-a", a.Bssid, iface)

	err := cmd.Start() // Do not wait

	cur_atk := attacks.Attack{
		Type:    "ArpReplay",
		Target:  a.Bssid,
		Running: false,
		Started: time.Now().String(),
	}

	if err != nil {
		cur_atk.Running = true
		cur_atk.Init(cmd.Process)
	}

	return cur_atk, err
}

var captures_nb = 0

// Start a capture process
func (a *AP) Capture(iface string) (attacks.Attack, string, error) {
	path := "go-wifi_capture-" + strconv.Itoa(captures_nb)
	captures_nb += 1

	// Make a specific dir so we do not mix captures
	// TODO: change mode
	err := os.Mkdir(path, 766)
	if err == nil {
		return nil, nil, err
	}

	path += "go-wifi"
	cmd := exec.Command("airodump-ng", "--write", path, "-c", a.Channel, "--output-format", "pcap", "--bssid", a.Bssid, iface)

	err = cmd.Start() // Do not wait

	cur_atk := attacks.Attack{
		Type:    "Capture",
		Target:  a.Bssid,
		Running: false,
		Started: time.Now().String(),
	}

	if err != nil {
		cur_atk.Running = true
		cur_atk.Init(cmd.Process)
	}

	// Because of an import cycle, we cannot build the Capture object, we just return the dir's path
	return cur_atk, path, err
}
