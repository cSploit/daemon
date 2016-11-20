package models

import "github.com/cSploit/daemon/models/internal"

//TODO: turn it into tcpdump capture, with a field which specify the physical medium type ( 802.11 or Ethernet )
//TODO: Handshake entity { nonce, hmac, ... }
//TODO: WpaKey entity { Ap, Handshake, Key }
//TODO: WepCrackJob { Capture, Handshake, Ap }

// TODO: make cracking and trying keys jobs

// an airodump capture file
type (
	Capture struct {
		internal.Base

		Key        *string `json:"key"`
		Handshake bool     `json:"has_handshake"`
		//Cracking   bool    `json:"cracking"`
		File       string  `json:"-"`

		Ap   AP   `json:"-"`
		ApId uint `json:"ap_id"`
	}

	Target struct {
		Bssid string `json:"bssid"`
		Essid string `json:"essid"`
		// WPA, WPA2, WEP, OPN
		Privacy string `json:"privacy"`
	}
)

// Build the struct thanks to the dir (with .pcap and .csv) path
func (c *Capture) Init(path_to_captures string) {
	c.File = path_to_captures + "go-wifi-01.cap"

	// Check if we have an Handshake
	if c.Ap.Privacy == "WPA" || c.Ap.Privacy == "WPA2" {
		c.checkForHandshake()
	}
}

// Return succesfull key
func (c *Capture) TryKeys(keys ...string) string {
	if c.Ap.Privacy == "WEP" || c.Ap.Privacy == "OPN" {
		// Only wpa
		return nil
	}

	// build a temp dict
	path := os.TempDir() + "go-wifi-tmp-dict"

	file, err := os.Create(path)
	if err != nil {
		// Got an error, exit
		return
	}
	defer file.Close()
	defer os.Remove(path)

	for _, key := range keys {
		file.WriteString(key + "\n")
	}

	return c.crackWPA(path)
}

// Return ascii key; if cracking WEP dict can be null
func (c *Capture) AttemptToCrack(dict string) string {
	// Do not crack a second time!
	if c.Key != nil {
		return c.Key
	}

	// Start here
	var key string

	if (c.Ap.Privacy == "WPA" || c.Ap.Privacy == "WPA2") && dict != nil {
		key = c.crackWPA(dict)
	} else if c.Ap.Privacy == "WEP" {
		key = c.crackWEP()
	} else {
		key = nil
	}

	if key != nil {
		c.Key = key
	}

	return key
}

func (c *Capture) crackWPA(dict string) string {
	// I use a random file so you can run the func in parallel
	path_to_key := os.TempDir() + "go-wifi_key" + strconv.Itoa(rand.Uint32())

	// If the file exist, delete it
	os.Remove(path_to_key)

	cmd := exec.Command("aircrack-ng", "-a", "2", "-l", path_to_key, "-w", dict, "-b", c.Ap.Bssid, c.File)
	cmd.Run()

	// Wait termination so we can get the key
	cmd.Wait()

	key_buf, err := ioutil.ReadFile(path_to_key)
	if err != nil {
		// no key found
		return nil
	}

	return string(key_buf)
}

func (c *Capture) crackWEP() string {
	// Start with PTW
	// I use a random file so you can run the func in parallel
	path_to_key := os.TempDir() + "go-wifi_key" + strconv.Itoa(rand.Uint32())

	// If the file exist, delete it
	os.Remove(path_to_key)

	cmd := exec.Command("aircrack-ng", "-D", "-z", "-a", "1", "-l", path_to_key, "-b", c.Ap.Bssid, c.File)
	cmd.Run()

	// Wait termination so we can get the key
	cmd.Wait()

	// Check if we succeed
	key_buf, err := ioutil.ReadFile(path_to_key)
	if err != nil {
		// no key found, start Korek
		cmd = exec.Command("aircrack-ng", "-D", "-K", "-a", "1", "-l", path_to_key, "-b", c.Ap.Bssid, c.File)
		cmd.Run()
		cmd.Wait()

		key_buf, err = ioutil.ReadFile(path_to_key)
		if err != nil {
			// Korek and PTW failed, exit
			return nil
		}
	}

	// key_buf has a key!
	return string(key_buf)
}

func (c *Capture) checkForHandshake() {
	// Thank you wifite (l. 2478, has_handshake_aircrack)
	// build a temp dict
	path := os.TempDir() + "fake-dict"

	file, err := os.Create(path)
	if err != nil {
		// Got an error, exit
		return
	}
	defer file.Close()

	file.WriteString("that_is_a_fake_key_no_one_will_use")

	cmd := exec.Command("aircrack-ng",  "-a", "2", "-w", path, "-b", c.Ap.Bssid, c.File)

	ouptut, err2 := cmd.Output()

	if err2 == nil {
		if strings.Contains(string(ouptut), "Passphrase not in dictionary") {
			c.Handshake = true
		} else {
			c.Handshake = false
		}
	}

	// Delete file
	os.Remove(path)
}
