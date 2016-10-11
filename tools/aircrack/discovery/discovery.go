package discovery

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/cSploit/daemon/tools/aircrack/AP"
)

// WARNING in order to use airodump-ng, you may need root access

// JSON exportable structs
// Full discovery
type Discovery struct {
	APs     []AP.AP     `json:"aps"`
	Clients []AP.Client `json:"clients"`
	Running bool        `json:"running"`
	Started string      `json:"started at"`
	Stopped string      `json:"stopped at"`
	process *os.Process
}

// Stop the discovery...
func (d *Discovery) Stop() error {
	err := d.process.Kill()
	if err != nil {
		return err
	}

	d.Running = false
	d.Stopped = time.Now().String()

	return nil
}

// THE most important function: parse /tmp/discovery-01.csv and fill the structs
// Might be nice to run as a goroutine...
// IDEA: parallelise the parsing with 2 goroutines
func (d *Discovery) Parse() error {
	// Dirty hack to have a clean dump
	dump, err := ioutil.ReadFile(os.TempDir() + "/discovery-01.csv")
	if err != nil {
		return err
	}

	dump_str := string(dump)
	// Replace endline with just an  \n
	dump_str = strings.Replace(dump_str, ", \r\n", ", \n", -1)
	dump_str = strings.Replace(dump_str, ",\r\n", ",\n", -1)
	dump_split := strings.SplitN(dump_str, "\r\n", 4)

	// Extract the two parts of the csv
	dump_aps := dump_split[2]
	dump_clients := dump_split[3]
	dump_clients = strings.SplitN(dump_clients, "\r\n", 2)[1]

	// End of dirty hack, fill the structs
	reader_aps := csv.NewReader(strings.NewReader(dump_aps))
	reader_clients := csv.NewReader(strings.NewReader(dump_clients))

	// We will fill them back later
	d.APs = nil
	d.Clients = nil

	// Start with the aps
	for {
		record, csv_err := reader_aps.Read()
		if csv_err == io.EOF {
			break
		}
		if csv_err != nil {
			return err
		}

		// Okay, fill an AP struct then append to the dump

		// TODO: clean that
		// NOTE: I am too lazy to check the errors
		channel, _ := strconv.Atoi(record[3])
		speed, _ := strconv.Atoi(record[4])
		power, _ := strconv.Atoi(record[8])
		beacons, _ := strconv.Atoi(record[9])
		ivs, _ := strconv.Atoi(record[10])
		idlen, _ := strconv.Atoi(record[12])

		cur_ap := AP.AP{
			Bssid:   record[0],
			First:   record[1],
			Last:    record[2],
			Channel: channel,
			Speed:   speed,
			Privacy: record[5],
			Cipher:  record[6],
			Auth:    record[7],
			Power:   power,
			Beacons: beacons,
			IVs:     ivs,
			Lan:     strings.Replace(record[11], " ", "", -1), // Clean blanks
			IdLen:   idlen,
			Essid:   record[13],
			Key:     record[14],
		}

		d.APs = append(d.APs, cur_ap)
	}

	// Continue with the clients
	for {
		record, csv_err := reader_clients.Read()
		if csv_err == io.EOF {
			break
		}
		if csv_err != nil {
			return err
		}

		// Okay, fill a Client struct then append to the dump

		// TODO: clean that
		// NOTE: too lazy to fix the errors
		power, _ := strconv.Atoi(record[3])
		packets, _ := strconv.Atoi(record[4])

		cur_client := AP.Client{
			Station: record[0],
			First:   record[1],
			Last:    record[2],
			Power:   power,
			Packets: packets,
			Bssid:   record[5],
			Probed:  record[6],
		}

		d.Clients = append(d.Clients, cur_client)
	}

	return nil
}

// Start a new discovery thanks to airodump-ng
// iface MUST be the name of valid monitor mode iface
// it will create a temp file named "discovery-01.csv", if it exist,
// it will be deleted!
// Return a Discovery object
func StartDiscovery(iface string) (Discovery, error) {
	// Delete previous log file
	os.Remove(os.TempDir() + "/discovery-01.csv")

	// okay, enough cosmetics, time for real code!
	cmd := exec.Command("airodump-ng", "--write", os.TempDir()+"/discovery", "--output-format", "csv", "--wps", iface)

	err := cmd.Start() // Do not wait

	discovery := Discovery{
		Started: time.Now().String(),
		Running: false,
	}

	if err == nil {
		discovery.Running = true
		discovery.process = cmd.Process
	}

	return discovery, err
}
