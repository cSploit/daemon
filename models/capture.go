package models

import (
	"github.com/cSploit/daemon/models/internal"
	"io/ioutil"
	"errors"
	"os"
	"strings"
	"strconv"
	"time"
)

//TODO: turn it into tcpdump capture, with a field which specify the physical medium type ( 802.11 or Ethernet )
//TODO: Handshake entity { nonce, hmac, ... }

// TODO: IVs

// TODO: trying keys jobs

// an airodump capture file
type Capture struct {
	internal.Base

	Key       string `json:"key"`
	Handshake bool    `json:"has_handshake"`
	Cracking  bool    `json:"cracking"`
	File      string  `json:"-"`

	Dict string `json:"dict"`

	Ap   AP   `json:"-"`
	ApId uint `json:"ap_id"`
}

var key_nb int

// Return ascii key; if cracking WEP dict can be null
func (c *Capture) Crack() (j Job, e error) {
	// Do not crack a second time!
	if c.Key != "" {
		e = errors.New("Already cracked")
		return
	}

	c.Cracking = true

	if (c.Ap.Privacy == "WPA" || c.Ap.Privacy == "WPA2") {
		if c.Dict != "" {
			j, e = c.crackWPA()
		} else {
			e = errors.New("Dictionnary needed for WPA(2) attack")
		}
	} else if c.Ap.Privacy == "WEP" {
		j, e = c.crackWEP()
	} else {
		e = errors.New("Target seems not to be encrypted")
	}

	return
}

func (c *Capture) crackWPA() (j Job, e error) {
	path_to_key := os.TempDir() + "go-wifi_key" + strconv.Itoa(key_nb)
	key_nb += 1

	pj, e := CreateProcessJob("aircrack-ng", "-a", "2", "-l", path_to_key, "-w", c.Dict, "-b", c.Ap.Bssid, c.File)

	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Name", "CrackWpa ["+c.Ap.Bssid+"]")
		db.Model(&j).Association("Aps").Append(c)
	}

	go c.waitCrack(pj, path_to_key)
	return
}

func (c *Capture) crackWEP() (j Job, e error) {
	path_to_key := os.TempDir() + "go-wifi_key" + strconv.Itoa(key_nb)
	key_nb += 1

	pj, e := CreateProcessJob("aircrack-ng", "-D", "-z", "-a", "1", "-l", path_to_key, "-b", c.Ap.Bssid, c.File)

	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Name", "CrackWep ["+c.Ap.Bssid+"]")
		db.Model(&j).Association("Aps").Append(c)
	}

	go c.waitCrack(pj, path_to_key)
	return
}

func (c *Capture) waitCrack(pj *ProcessJob, path_to_key string) {
	for {
		if pj.ExitStatus == nil {
			time.Sleep(time.Second * 1)
		}
	}

	key_buff, err := ioutil.ReadFile(path_to_key)
	if err == nil {
		c.Key = string(key_buff)
	}

	c.Cracking = false
}

func (c *Capture) CheckForHandshake() (j Job, e error) {
	// Thank you wifite (l. 2478, has_handshake_aircrack)
	// build a temp dict
	path := os.TempDir() + "fake-dict"

	file, e := os.Create(path)
	if e != nil {
		// Got an error, exit
		return
	}
	defer file.Close()

	file.WriteString("that_is_a_fake_key_no_one_will_use")

	pj, e := CreateProcessJob("aircrack-ng", "-a", "2", "-w", path, "-b", c.Ap.Bssid, c.File)

	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Name", "CheckHandshake ["+c.Ap.Bssid+"]")
		db.Model(&j).Association("Aps").Append(c)
	}

	go c.waitHandshakeTester(pj, file)
	return
}

func (c *Capture) waitHandshakeTester(pj *ProcessJob, file *os.File) {
	for {
		if pj.ExitStatus == nil {
			time.Sleep(time.Second * 1)
		}
	}

	if strings.Contains(pj.Output, "Passphrase not in dictionary") {
		c.Handshake = true
	} else {
		c.Handshake = false
	}

	os.Remove(file.Name())
}
