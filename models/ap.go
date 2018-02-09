package models

import (
	"github.com/cSploit/daemon/models/internal"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	internal.RegisterModels(&AP{})
}

// Access Point ( courtesy of aircrack )
type AP struct {
	internal.Base
	Bssid   string    `json:"bssid"`
	First   time.Time `json:"first_seen"`
	Last    time.Time `json:"last_seen"`
	Channel int       `json:"channel"`
	Speed   int       `json:"speed"`
	Privacy string    `json:"privacy"`
	Cipher  string    `json:"cipher"`
	Auth    string    `json:"auth"`
	Power   int       `json:"power"`
	Beacons int       `json:"beacons"`
	IVs     int       `json:"ivs"`
	Lan     string    `json:"lan_ip"`
	IdLen   int       `json:"id_len"`
	Essid   string    `json:"essid"`
	Key     string    `json:"key"`
	//Wps     bool   `json:"wps"`

	// Does the fake auth succeed?
	FakeAuthed bool `json:"fake_auth"`

	Iface   Iface `json:"-"`
	IfaceId uint  `json:"-"`
	Jobs    []Job `json:"-" gorm:"many2many:job_aps;"`
}

// DEAUTH infinitely the AP using broadcast address
func (a *AP) Deauth() (j Job, e error) {
	pj, e := CreateProcessJob("aireplay-ng", "-0", "0", "-a", a.Bssid, a.Iface.Name)

	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Name", "Deauth ["+a.Bssid+"]")
		db.Model(&j).Association("Aps").Append(a)
	}

	return
}

// Try a fake auth on the ap
func (a *AP) FakeAuth() (j Job, e error) {
	pj, e := CreateProcessJob("aireplay-ng", "-1", "0", "-a", a.Bssid, "-T", "1", a.Iface.Name)

	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Name", "FakeAuth ["+a.Bssid+"]")
		db.Model(&j).Association("Aps").Append(a)
	}

	go a.checkFakeAuth(pj)

	return
}

func (a *AP) checkFakeAuth(pj *ProcessJob) {
	for {
		if pj.ExitStatus == nil {
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}

	if strings.Contains(pj.Output, "Association successful") {
		a.FakeAuthed = true
	} else {
		a.FakeAuthed = false
	}
}

// ARP replay!!
func (a *AP) ArpReplay(iface string) (j Job, e error) {
	pj, e := CreateProcessJob("aireplay-ng", "-3", "-a", a.Bssid, a.Iface.Name)

	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Name", "ArpReplay ["+a.Bssid+"]")
		db.Model(&j).Association("Aps").Append(a)
	}

	return
}

var captures_nb = 0

// Start a capture process
func (a *AP) Capture() (j Job, e error) {
	path := "go-wifi_capture-" + strconv.Itoa(captures_nb)
	captures_nb += 1

	// Make a specific dir so we do not mix captures
	err := os.Mkdir(path, 0755)
	if err != nil {
		log.Error(err)
	}

	path += "/go-wifi"
	ch := strconv.Itoa(a.Channel)
	pj, e := CreateProcessJob("airodump-ng", "--write", path, "-c", ch, "--output-format", "pcap", "--bssid", a.Bssid, a.Iface.Name)

	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Name", "Capture ["+a.Bssid+"]")
		db.Model(&j).Association("Aps").Append(a)

		//TODO: start a routine that update the Capture record
		capture := &Capture{Ap: *a, ApId: a.ID, File: path + "-01.pcap"}
		db.Save(capture)
	}

	return
}

func FindAp(id uint) (a *AP, e error) {
	a = &AP{}
	e = internal.Db.Find(a, id).Error
	return
}

func FindApByBssid(bssid string) (a *AP, e error) {
	a = &AP{}
	e = internal.Db.Where("bssid = ?", bssid).Find(a).Error
	return
}
