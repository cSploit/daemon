package models

import (
	"github.com/cSploit/daemon/models/internal"
	"os"
	"strconv"
	"errors"
)

// TODO: add missing hostapd options (at least for wpa, wep, wps and hidden aps)

func init() {
	internal.RegisterModels(&RogueAP{})
}

type (
	HostapdConf struct {
		// Attract other clients
		// enable_mana=1
		EnableMana bool `json:"enable_mana"`
		// By default, MANA will be a little stealthy and only advertise probed for networks
		// directly to the device that probed for it.
		// However, not all devices probe as much as they used to, and some devices will
		// probe with "random" locally administered MAC addresses.
		// Loud mode will re-broadcast all networks to all devices.
		// mana_loud=0
		ManaLoud bool `json:"mana_loud"`
		// Normal access points MAC ACLs will only work at association level. This option
		// will expand MAC ACLs to probe responses.
		// It requires macaddr_acl to be set later in the config file to work. This controls
		// whether we're operating in black or white list mode. The MACs are defined in the
		// files listed in accept_mac_file and deny_mac_file.
		// Setting ignore_broadcast_ssid below will also hide the base network from
		// non-authorised devices.
		// mana_macacl=0
		ManaMacAcl bool `json:"mana_mac_acl"`
		// hostap/wired/madwifi/test/none/nl80211/bsd
		// nl80211 is fine
		// driver=nl80211
		Driver string `json:"driver"`
		// ssid=Hack
		SSID string `json:"ssid"`
		// bssid=02:21:91:01:11:31
		BSSID string `json:"bssid"`
		// channel=6
		Channel int `json:"channel"`
		// 0: accept unless in deny list
		// 1: deny unless in accept list
		// macaddr_acl=0
		MacAddrAcl int `json:"mac_addr_acl"`
		// accept_mac_file=/etc/mana-toolkit/hostapd.accept
		AcceptMacFile string `json:"accept_mac_file"`
		// deny_mac_file=/etc/mana-toolkit/hostapd.deny
		DenyMacFile string `json:"deny_mac_file"`
		// auth_algs=3
		AuthAlgs int `json:"auth_algs"`
	}

	RogueAP struct {
		Hostapd HostapdConf `json:"hostapd"`
		Iface Iface `json:"-"`
	}
)

func (r *RogueAP) StartHostapd() (j Job, e error) {
	// We first build the conf file
	path := os.TempDir() + "hostapd.conf"
	// Remove previous version
	os.Remove(path)

	// Start a new conf
	file, e := os.Create(path)
	if e != nil {
		return
	}
	defer file.Close()

	// Explicitly enable or disable options

	// Handle mana options
	if r.Hostapd.EnableMana {
		file.WriteString("enable_mana=1\n")

		if r.Hostapd.ManaLoud {
			file.WriteString("mana_loud=1\n")
		} else {
			file.WriteString("mana_loud=0\n")
		}

		if r.Hostapd.ManaMacAcl {
			file.WriteString("mana_macacl=1\n")
		} else {
			file.WriteString("mana_macacl=0\n")
		}
	} else {
		file.WriteString("enable_mana=0\n")
	}

	// TODO: check if the interface is a valid one
	if r.Iface.Name != "" {
		file.WriteString("interface=" + r.Iface.Name + "\n")
	} else {
		e = errors.New("Need an interface")
		return
	}

	if r.Hostapd.SSID != "" {
		file.WriteString("ssid=" + r.Hostapd.SSID + "\n")
	} else {
		e = errors.New("Need a SSID")
		return
	}

	// TODO: check if it's a valid one
	if r.Hostapd.BSSID != nil {
		file.WriteString("bssid=" + r.Hostapd.BSSID + "\n")
	}

	// TODO: check if it's a valid one
	if r.Hostapd.Channel != nil {
		file.WriteString("channel=" + strconv.Itoa(r.Hostapd.Channel) + "\n")
	}

	if r.Hostapd.MacAddrAcl == 0 {
		file.WriteString("macaddr_acl=0\n")
	} else if r.Hostapd.MacAddrAcl == 1 {
		file.WriteString("macaddr_acl=1\n")
	} else {
		e = errors.New("Invalid MacAddrAcl value, should be 0 or 1")
		return
	}

	// TODO: check if files exists
	if r.Hostapd.AcceptMacFile != "" {
		file.WriteString("accept_mac_file=" + r.Hostapd.AcceptMacFile + "\n")
	}
	if r.Hostapd.DenyMacFile != "" {
		file.WriteString("deny_mac_file=" + r.Hostapd.DenyMacFile + "\n")
	}

	// Do not support an other value ATM
	if r.Hostapd.AuthAlgs == 3 {
		file.WriteString("auth_algs=3\n")
	} else {
		e = errors.New("AuthAlgs must be 3")
		return
	}

	// Time to start hostapd!
	// Assume the hostapd tool is patched and renamed to prevent using the standard version
	pj, e := CreateProcessJob("hostapd-mana", path)

	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Mana rogue AP [hostapd]")
		db.Model(&j).Association("rogue_ap").Append(r)
	}

	return
}

func (r *RogueAP) SetupIface() (j Job, e error) {
	pj, e := CreateProcessJob("ifconfig", r.Iface.Name, "10.0.0.1", "netmask", "255.255.255.0")

	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Mana rogue AP [ifconfig]")
		db.Model(&j).Association("rogue_ap").Append(r)
	}

	return
}

func (r *RogueAP) SetupRoute() (j Job, error) {
	pj, e := CreateProcessJob("route", "add", "-net", "10.0.0.0", "netmask", "255.255.255.0", "gw", "10.0.0.1")

	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Mana rogue AP [route]")
		db.Model(&j).Association("rogue_ap").Append(r)
	}

	return
}

func (r *RogueAP) StartDnsMasq() (j Job, e error) {
	// We first build the conf file
	path := os.TempDir() + "dnsmasq.conf"
	// Remove previous version
	os.Remove(path)

	// Start a new conf
	file, e := os.Create(path)
	if e != nil {
		return
	}
	defer file.Close()

	// Hardcoded
	file.WriteString("dhcp-range=10.0.0.100,10.0.0.254,1h\n")
	// DNS
	file.WriteString("dhcp-option=6,8.8.8.8\n")
	// Gateway
	file.WriteString("dhcp-option=3,10.0.0.1\n")
	file.WriteString("dhcp-authoritative\n")

	pj, e := CreateProcessJob("dnsmasq", "-z", "-C", path, "-i", r.Iface.Name, "-I", "lo", "-k")


	if e == nil {
		j = pj.Job
		db := internal.Db
		db.Model(&j).Update("Mana rogue AP [dnsmasq]")
		db.Model(&j).Association("rogue_ap").Append(r)
	}

	return
}

func FindRogueAp(id uint) (r *RogueAP, e error) {
	r = &RogueAP{}
	e = internal.Db.Find(r, id).Error
	return
}