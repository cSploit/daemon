package captures

import (
	"../AP"
)

// Capture struct: handle airodump captures to crack them with aircrack-ng
type Capture struct {
	Key       string `json:"key"`
	Target    AP.AP  `json:"target"`
	Handshake bool   `json:"handshake captured"`
	IVs       int    `json:"ivs"`
	Cracking  bool   `json:"trying to crack"`
	pcap_file string
}

// Build the struct thanks to the dir (with .pcap and .csv) path
func (c *Capture) Init(path_to_captures string, target AP.AP) {
	c.pcap_file = path_to_captures + "go-wifi-01.cap"

	// Fill the struct!
	c.Target = target // Everything is there (privacy etc...)

	// Check if we have an Handshake
	if target.Privacy == "WPA" || target.Privacy == "WPA2" {
		c.checkForHandshake()
	}

	c.getIVs()

}

// Return succesfull key
func (c *Capture) TryKeys(...string) string {

}

// Return success, ascii key
func (c *Capture) AttemptToCrack() (bool, string) {

}

func (c *Capture) checkForHandshake() {

}

func (c *Capture) getIVs() {

}
