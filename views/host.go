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
package views

import "github.com/cSploit/daemon/models"

type hostIdxElem struct {
	models.Host
	OpenPortCount int    `json:"open_port_count"`
	HidePorts     string `json:"ports,omitempty"`
}

type hostShowView struct {
	models.Host
	HideNetworkID string      `json:"network_id,omitempty"`
	PortsView     interface{} `json:"ports"`
	NetworkView   interface{} `json:"network,omitempty"`
}

func HostsIndex(arg interface{}) interface{} {
	hosts := arg.([]models.Host)
	res := make([]hostIdxElem, len(hosts))

	for i, h := range hosts {
		// we assume that h.Ports contains all
		// and only the open ports

		res[i] = hostIdxElem{
			Host:          h,
			OpenPortCount: len(h.Ports),
		}
	}

	return res
}

func HostsShow(arg interface{}) interface{} {
	host := arg.(models.Host)
	var net interface{}

	portsView := PortIndex(host.Ports)

	if host.Network != nil {
		net = networkAsChild(*host.Network)
	}

	res := hostShowView{
		Host:        host,
		PortsView:   portsView,
		NetworkView: net,
	}

	return res
}
