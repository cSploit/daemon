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
	netHelper "github.com/cSploit/daemon/helpers/net"
	"github.com/cSploit/daemon/models/internal"
	"net"
)

func init() {
	internal.RegisterModels(&Network{})
}

type Network struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	IfaceName string `json:"iface_name"`
	IpAddr    string `json:"ip_addr"`
	Hosts     []Host `json:"hosts"`
}

func NewNetwork(ifName, ipAddr string) *Network {
	return &Network{
		IfaceName: ifName,
		IpAddr:    ipAddr,
	}
}

func FindNetwork(ipNet *net.IPNet) *Network {
	network := &Network{}

	dbRes := internal.Db.Where("ip_addr = ?", ipNet.String()).Find(network)

	if dbRes.RecordNotFound() {
		return nil
	} else if dbRes.Error != nil {
		log.Warning(dbRes.Error)
		return nil
	}

	return network
}

func CreateNetwork(ipNet *net.IPNet) *Network {
	var ifName string

	if iface, err := netHelper.InterfaceForIp(ipNet.IP); err != nil {
		log.Error(err)
		ifName = "unknown"
	} else {
		ifName = iface.Name
	}

	network := NewNetwork(ifName, ipNet.String())

	dbRes := internal.Db.Create(network)

	if dbRes.Error != nil {
		log.Error(dbRes.Error)
		return nil
	}

	return network
}

func FindOrCreateNetwork(ipNet *net.IPNet) *Network {
	res := FindNetwork(ipNet)

	if res == nil {
		res = CreateNetwork(ipNet)
	}

	return res
}

func (n *Network) GetHosts() []Host {
	var hosts []Host

	dbRes := internal.Db.Where("network_id = ?", n.ID).Find(&hosts)

	if dbRes.Error != nil {
		log.Error(dbRes.Error)
		return hosts
	}

	return hosts
}
