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

type networkIdxElem struct {
	models.Network
	HideHosts string `json:"hosts,omitempty"`
}

type networkShowView struct {
	models.Network
	OverrideHosts interface{} `json:"hosts,omitempty"`
}

func NetworkIndex(args interface{}) interface{} {
	nets := args.([]models.Network)
	res := make([]networkIdxElem, len(nets))

	for i, n := range nets {
		res[i] = networkIdxElem{Network: n}
	}

	return res
}

func NetworkShow(arg interface{}) interface{} {
	net := arg.(models.Network)
	res := networkShowView{Network: net}

	if len(net.Hosts) > 0 {
		res.OverrideHosts = HostsIndex(net.Hosts)
	}

	return res
}

func networkAsChild(arg interface{}) interface{} {
	network := arg.(models.Network)
	return networkIdxElem{Network: network}
}
