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
package models

import (
	"github.com/lair-framework/go-nmap"
)

type Port struct {
	ID       uint     `json:"id"`
	HostId   uint     `json:"host_id"`
	Protocol string   `json:"protocol"` // (ip|tcp|udp|sctp)
	Number   int      `json:"number"`
	State    string   `json:"state"` // "open","filtered","unfiltered","closed","open|filtered","closed|filtered","unknown"
	Service  *Service `json:"-"`
}

func NewPort(p nmap.Port) *Port {

	res := &Port{Protocol: p.Protocol, Number: p.PortId, State: p.State.State}

	if p.Service.Name != "" && p.Service.Name != "unknown" {
		res.Service = NewService(p.Service)
	}

	return res
}
