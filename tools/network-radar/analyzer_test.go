/* cSploit - a simple penetration testing suite
 * Copyright (C) 2016 Massimo Dragano aka tux_mind <tux_mind@csploit.org>
 *
 * cSploit is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * cSploit is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with cSploit.  If not, see <http://www.gnu.org/licenses/\>.
 *
 */
package network_radar

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"testing"
)

// took from wireshark
// a good NetBIOS query to my router, ETH + IP + UDP + NetBIOS

var goodNBQueryPkt = []byte{
	0x00, 0x26, 0x5a, 0x9d, 0xf0, 0x76, 0x64, 0x70, 0x02, 0xda, 0x03, 0x05, 0x08, 0x00, 0x45, 0x00,
	0x00, 0x4e, 0x08, 0x8b, 0x40, 0x00, 0x40, 0x11, 0xb0, 0xab, 0xc0, 0xa8, 0x00, 0x17, 0xc0, 0xa8,
	0x00, 0x01, 0xb0, 0x49, 0x00, 0x89, 0x00, 0x3a, 0x0c, 0xdd, 0x82, 0x28, 0x00, 0x00, 0x00, 0x01,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x43, 0x4b, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
	0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
	0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x00, 0x00, 0x21, 0x00, 0x01,
}

var arpRequestPkt = []byte{
	0, 4, 0, 1, 0, 6, 100, 112, 2, 218, 3, 5, 0, 0, 8, 6, // Linux cooked capture
	0, 0, 4, 0, 1, 0, 6, 100, 112, 2, 218, 3, 5, 0, 0, 8, 6, 0, 0, 4, 0, 1, 0, 6, 100, 112, 2, 218, 3, 5, 0, 0, 8, 6, 0,
}

func TestAnalyzeNetBIOS(t *testing.T) {
	pkt := gopacket.NewPacket(goodNBQueryPkt, layers.LayerTypeEthernet, gopacket.Default)

	for _, l := range pkt.Layers() {
		t.Logf("contains: %v", l.LayerType())
	}

	if pkt.NetworkLayer() == nil {
		t.Error("created packet does not implement Network layer")
		t.Fail()
	}
}

func TestAnalyzeARP(t *testing.T) {
	pkt := gopacket.NewPacket(arpRequestPkt, layers.LayerTypeLinuxSLL, gopacket.Default)

	for _, l := range pkt.Layers() {
		t.Logf("contains: %v", l.LayerType())
	}
}
