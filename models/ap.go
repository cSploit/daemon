package models

import "github.com/cSploit/daemon/models/internal"

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
	//TODO: discovered_on []Iface -> use this to known if this AP is reachable by multiple adapters
	Jobs []Job `json:"-" gorm:"many2many:job_aps;"`
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
