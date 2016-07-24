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
package model

/*

module that interact with cSploit model to fetch hosts
and update them as needed

*/

import (
	netHelper "github.com/cSploit/daemon/helpers/net"
	"github.com/cSploit/daemon/models"
	"github.com/op/go-logging"
	"golang.org/x/net/context"

	"gopkg.in/guregu/null.v3"
	"net"
	"time"
)

var log = logging.MustGetLogger("netowrk_radar.model")

func FindNetwork(ipNet *net.IPNet) *models.Network {
	db := models.GetDbInstance()
	network := &models.Network{}

	dbRes := db.Where("ip_addr = ?", ipNet.String()).Find(network)

	if dbRes.RecordNotFound() {
		return nil
	} else if dbRes.Error != nil {
		log.Warning(dbRes.Error)
		return nil
	}

	return network
}

func CreateNetwork(ipNet *net.IPNet) *models.Network {
	var ifName string

	if iface, err := netHelper.InterfaceForIp(ipNet.IP); err != nil {
		log.Error(err)
		ifName = "unknown"
	} else {
		ifName = iface.Name
	}

	db := models.GetDbInstance()

	network := models.NewNetwork(ifName, ipNet.String())

	dbRes := db.Create(network)

	if dbRes.Error != nil {
		log.Error(dbRes.Error)
		return nil
	}

	return network
}

func walkHostsIPByNetID(ctx context.Context, netID uint) <-chan net.IP {
	c := make(chan net.IP)

	go func() {
		var hosts []models.Host

		defer close(c)

		db := models.GetDbInstance()
		dbRes := db.Where("network_id = ?", netID).Find(&hosts)

		if dbRes.Error != nil {
			log.Error(dbRes.Error)
			return
		}

		for _, h := range hosts {
			ip := net.ParseIP(h.IpAddr)

			if ip == nil {
				log.Warningf("unable to parse ip '%s' for host %s", h.IpAddr, h)
				continue
			}

			select {
			case c <- ip:
			case <-ctx.Done():
				return
			}
		}
	}()

	return c
}

func NewKnownHostsWalker(ipNet *net.IPNet) KnownHostsIPWalker {
	found := FindNetwork(ipNet)

	if found == nil {
		found = CreateNetwork(ipNet)
		if found == nil {
			return nil
		}
	}

	netId := found.ID

	return func(ctx context.Context) <-chan net.IP {
		return walkHostsIPByNetID(ctx, netId)
	}
}

func NotifyHostSeen(hwAddr net.HardwareAddr, ipAddr net.IP, name string) error {
	hwId, err := netHelper.MacAddrToUInt(hwAddr)

	if err != nil {
		return err
	}

	var HwAddrEntity models.HwAddr

	db := models.GetDbInstance()

	dbRes := db.Preload("Host").Find(&HwAddrEntity, hwId)

	if dbRes.RecordNotFound() {
		return onNewHost(hwAddr, ipAddr, name)
	} else if dbRes.Error != nil {
		return dbRes.Error
	}

	host := HwAddrEntity.Host

	if host == nil {
		return onNewHostWithHwAddr(&HwAddrEntity, ipAddr, name)
	}

	return onHostSeen(host, ipAddr, name)
}

//TODO: fire an event for each of these functions

func onNewHost(hwAddr net.HardwareAddr, ipAddr net.IP, name string) error {
	hw, err := models.NewHwAddr(hwAddr)

	if err != nil {
		return err
	}

	return onNewHostWithHwAddr(hw, ipAddr, name)
}

func onNewHostWithHwAddr(hwAddr *models.HwAddr, ipAddr net.IP, name string) error {
	db := models.GetDbInstance()

	nullName := null.NewString(name, len(name) > 0)

	host := models.Host{HwAddr: hwAddr, IpAddr: ipAddr.String(), Name: nullName}

	return db.Create(&host).Error
}

func onHostSeen(host *models.Host, ipAddr net.IP, name string) error {
	db := models.GetDbInstance()
	host.IpAddr = ipAddr.String()
	host.Name = null.NewString(name, len(name) > 0)
	host.UpdatedAt = time.Now()
	return db.Save(host).Error
}
