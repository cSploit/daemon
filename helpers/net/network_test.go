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
package net

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"github.com/lair-framework/go-nmap"
)

func TestBuildBroadcastAddress(t *testing.T) {
	_, ipNet, _ := net.ParseCIDR("192.168.0.1/24")
	brAddr := net.ParseIP("192.168.0.255")

	res := BuildBroadcastAddress(ipNet)

	assert.True(t, brAddr.Equal(res))
}

func TestNextIP(t *testing.T) {
	cur := net.ParseIP("192.168.1.255")
	next := net.ParseIP("192.168.2.0")
	NextIP(cur)

	assert.True(t, next.Equal(cur))
}

func TestNumHosts(t *testing.T) {
	_, ipNet, _ := net.ParseCIDR("192.168.0.1/27")
	a := assert.New(t)

	res := NumHosts(ipNet)

	a.Equal(uint64(30), res)

	_, ipNet, _ = net.ParseCIDR("192.168.0.1/16")
	res = NumHosts(ipNet)

	a.Equal(uint64(65534), res)

	_, ipNet, _ = net.ParseCIDR("192.168.0.1/0")
	res = NumHosts(ipNet)

	a.Equal(uint64(4294967296)-2, res)

	_, ipNet, _ = net.ParseCIDR("127.0.0.1/8")
	ipNet.IP = net.IPv4(127, 0, 0, 1)
	res = NumHosts(ipNet)

	a.Equal(uint64(16777216)-2, res)
}

func TestMacConversion(t *testing.T) {
	mac1, _ := net.ParseMAC("01:23:45:67:89:ab")
	mac2, _ := net.ParseMAC("01:23:45:67:89:ab:cd:ef")
	mac3, _ := net.ParseMAC("01:23:45:67:89:ab:cd:ef:00:00:01:23:45:67:89:ab:cd:ef:00:00")

	if _, err := MacAddrToUInt(mac1); err != nil {
		t.Fatalf("unable to convert %v to id: %v", mac1, err)
	}
	if _, err := MacAddrToUInt(mac2); err != nil {
		t.Fatalf("unable to convert %v to id: %v", mac2, err)
	}
	if _, err := MacAddrToUInt(mac3); err != nil {
		t.Fatalf("unable to convert %v to id: %v", mac3, err)
	}
}

func TestMACStringToUInt(t *testing.T) {
	samples := []struct {
		Addr string
		Val  uint64
	}{
		{"68:a3:c4:6f:fb:88", 115052584631176},
	}

	for _, s := range samples {
		m, err := net.ParseMAC(s.Addr)

		if err != nil {
			t.Fatalf("Sample MAC '%s' is broken, please fix it: %v", s.Addr, err)
		}

		n := nmap.Address{Addr: s.Addr, AddrType: "mac", Vendor: "Cisco"}

		for _, i := range []interface{}{m, n, s.Addr} {
			res, err := ParseHwAddr(i)

			if err != nil {
				t.Errorf("failed to create HwAddr from interface %T: %v", i, err)
			} else if res != s.Val {
				t.Errorf("using interface %T: expected %v, got %v", i, s.Val, res)
			}
		}
	}
}