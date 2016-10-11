package captures

import (
	"os"

	"github.com/cSploit/daemon/tools/aircrack/AP"
)

// Capture struct: handle airodump captures to crack them with aircrack-ng
type Capture struct {
	File      string `json:"file"`
	Key       string `json:"key"`
	Target    AP.AP  `json:"target"`
	Handshake bool   `json:"handshake_captured"`
	IVs       int    `json:"ivs"`
	Pkts      int    `json:"packets"`
	Cracking  bool   `json:"trying_to_crack"`
	process   *os.Process
}
