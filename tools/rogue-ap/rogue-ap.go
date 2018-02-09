package rogue_ap

import (
	"errors"
	"github.com/cSploit/daemon/models"
	"github.com/op/go-logging"
	"golang.org/x/net/context"
	"net"
	"os"
	"strings"
)

var (
	log = logging.MustGetLogger("rogue-ap")
)

type (
	RogueAP struct {
		Mana     bool
		Loud     bool
		DenyMac  []string
		AllowMac []string
		SSID     string
		BSSID    string
		Channel  int

		Iface models.Iface

		ctx    context.Context
		cancel context.CancelFunc
	}
)

func (r *RogueAP) Start() error {
	// Build the context
	r.ctx, r.cancel = context.WithCancel(context.Background())

	// Check some options

	// the net pkg parse lower case addresses
	if r.BSSID != "" {
		r.BSSID = strings.ToLower(r.BSSID)

		_, e := net.ParseMAC(r.BSSID)
		if e != nil {
			log.Error(e)
			return e
		}
	}

	if r.DenyMac != nil {
		for _, mac := range r.DenyMac {
			_, e := net.ParseMAC(strings.ToLower(mac))
			if e != nil {
				log.Error(e)
				return e
			}
		}
	}

	if r.AllowMac != nil {
		for _, mac := range r.AllowMac {
			_, e := net.ParseMAC(strings.ToLower(mac))
			if e != nil {
				log.Error(e)
				return e
			}
		}
	}

	if r.Iface.Name == "" {
		log.Error("No interface name")
		return errors.New("No iface defined")
	}

	// Done, configure
	log.Debug("Building config files...")

	path := os.TempDir()

	// Delete previous files
	os.Remove(path + "/hostapd.conf")
	os.Remove(path + "/dnsmasq.conf")

	// Start with hostapd
	file, e := os.Create(path + "/hostapd.conf")
	if e != nil {
		log.Error(e)
		return e
	}
	defer file.Close()

	if r.Mana {
		file.WriteString("enable_mana=1\n")

		if r.Loud {
			file.WriteString("mana_loud=1\n")
		}

		if r.DenyMac != nil || r.AllowMac != nil {
			file.WriteString("mana_macacl=1\n")
		}
	}

	if r.DenyMac != nil {
		file.WriteString("macaddr_acl=0\n")

		// Create blacklist
		os.Remove(path + "/hostapd.deny")

		deny, e := os.Create(path + "/hostapd.deny")
		if e != nil {
			log.Error(e)
			return e
		}
		defer deny.Close()

		for _, mac := range r.DenyMac {
			deny.WriteString(strings.ToLower(mac) + "\n")
		}

		// Finally, configure it
		file.WriteString("deny_mac_file=" + path + "/hostapd.deny\n")
	} else if r.AllowMac != nil {
		file.WriteString("macaddr_acl=1\n")

		// Create whitelist
		os.Remove(path + "/hostapd.accept")

		accept, e := os.Create(path + "/hostapd.accept")
		if e != nil {
			log.Error(e)
			return e
		}
		defer accept.Close()

		for _, mac := range r.AllowMac {
			accept.WriteString(strings.ToLower(mac) + "\n")
		}

		// Finally, configure it
		file.WriteString("accept_mac_file=" + path + "/hostapd.accept\n")
	}

	if r.SSID != "" {
		file.WriteString("ssid=" + r.SSID + "\n")
	} else {
		log.Debug("No ssid defined")
	}

	if r.BSSID != "" {
		file.WriteString("bssid=" + r.BSSID + "\n")
	} else {
		log.Debug("No bssid defined, will use default one")
	}

	file.WriteString("interface=" + r.Iface.Name + "\n")
	file.WriteString("auth_algs=3\n")

	log.Debug("Hostapd configured")

	// Continue with dnsmasq
	dnsmasq, e := os.Create(path + "/dnsmasq.conf")
	if e != nil {
		log.Error(e)
		return e
	}
	defer dnsmasq.Close()

	dnsmasq.WriteString("dhcp-range=10.0.0.100,10.0.0.254,1h\n")
	// DNS
	dnsmasq.WriteString("dhcp-option=6,8.8.8.8\n")
	// Gateway
	dnsmasq.WriteString("dhcp-option=3,10.0.0.1\n")
	dnsmasq.WriteString("dhcp-authoritative\n")

	log.Debug("Dnsmasq configured")

	// We will need to wait for some process, so we run it in a goroutine
	go r.startProcesses(path)

	return nil
}

func (r *RogueAP) startProcesses(path string) {
	// Start everything
	hostapd, e1 := models.CreateProcessJob("hostapd-mana", path+"/hostapd.conf")
	ifconfig, e2 := models.CreateProcessJob("ifconfig", r.Iface.Name, "10.0.0.1", "netmask", "255.255.255.0")
	route, e3 := models.CreateProcessJob("route", "add", "-net", "10.0.0.0", "netmask", "255.255.255.0", "gw", "10.0.0.1")
	dnsmasq, e4 := models.CreateProcessJob("dnsmasq", "-z", "-C", path+"/dnsmasq.conf", "-i", r.Iface.Name, "-I", "lo", "-k")

	if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
		log.Error("Got an error while starting processes")
		r.cancel()
		return
	}

	// We use the context to know if we need to stop...
	// Might be a bit ugly...
	for {
		select {
		// TODO: if hostapd or dnsmasq completed, kill everything
		case <-r.ctx.Done():
			// Kill only if we need to
			hostapd.Kill()
			ifconfig.Kill()
			route.Kill()
			dnsmasq.Kill()
		}
	}
}
