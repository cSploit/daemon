package models

import "github.com/cSploit/daemon/models/internal"

//TODO: turn it into tcpdump capture, with a field which specify the physical medium type ( 802.11 or Ethernet )
//TODO: Handshake entity { nonce, hmac, ... }
//TODO: WpaKey entity { Ap, Handshake, Key }
//TODO: WepCrackJob { Capture, Handshake, Ap }

// an airodump capture file
type Capture struct {
	internal.Base

	Key        *string `json:"key"`
	Handshakes int     `json:"handshakes"`
	IVs        int     `json:"ivs"`
	Cracking   bool    `json:"cracking"`
	File       string  `json:"-"`

	Ap   AP   `json:"-"`
	ApId uint `json:"ap_id"`
}
