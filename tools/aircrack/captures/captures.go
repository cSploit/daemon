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
	Handshake bool   `json:"handshake captured"`
	IVs       int    `json:"ivs"`
	Pkts      int    `json:"packets"`
	Cracking  bool   `json:"trying to crack"`
	process   *os.Process
}
