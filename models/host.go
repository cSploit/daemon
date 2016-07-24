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
	"github.com/op/go-logging"
	"gopkg.in/guregu/null.v3"
	"time"
)

var log = logging.MustGetLogger("daemon")

type Host struct {
	ID        uint        `json:"id"`
	CreatedAt time.Time   `json:"first_seen"`
	UpdatedAt time.Time   `json:"last_seen"`
	Name      null.String `json:"name"`
	IpAddr    string      `gorm:"index" json:"ip_addr"`
	HwAddr    *HwAddr     `json:"hw_addr"`
	Ports     []Port      `json:"ports"`
	Network   *Network    `json:"-"`
	NetworkID uint        `json:"network_id,omitempty"`
}

func NewHost(h nmap.Host) *Host {
	res := new(Host)

	res.Ports = make([]Port, 0)

	for _, p := range h.Ports {
		res.Ports = append(res.Ports, *NewPort(p))
	}

	for _, a := range h.Addresses {
		if a.AddrType == "mac" {
			var err error
			res.HwAddr, err = NewHwAddr(a)

			if err != nil {
				log.Warningf("unable to load MAC address: %v", err)
			}

			log.Debugf("created HW Addr: %v", res.HwAddr)
		} else {
			res.IpAddr = a.Addr
		}
	}

	return res
}
