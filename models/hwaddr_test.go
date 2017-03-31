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
package models_test

import (
	"github.com/cSploit/daemon/models"
	"github.com/lair-framework/go-nmap"
	"net"
	"testing"
)

var samples = []struct {
	Addr string
	Val  uint64
}{
	{"68:a3:c4:6f:fb:88", 115052584631176},
}

func TestNewHwAddr(t *testing.T) {
	for _, s := range samples {
		m, err := net.ParseMAC(s.Addr)

		if err != nil {
			t.Errorf("Sample MAC '%s' is broken, please fix it: %v", s.Addr, err)
		}

		n := nmap.Address{Addr: s.Addr, AddrType: "mac", Vendor: "Cisco"}

		for _, i := range []interface{}{m, n, s.Addr} {
			res, err := models.NewHwAddr(i)

			if err != nil {
				t.Errorf("failed to create HwAddr from interface %T: %v", i, err)
			}

			if res.ID != s.Val {
				t.Errorf("using interface %T: expected %v, got %v", i, s.Val, res.ID)
			}
		}

		res, err := models.NewHwAddr(n)

		if res.Vendor != n.Vendor {
			t.Errorf("expected vendor %s, got %s", n.Vendor, res.Vendor)
		}
	}
}
