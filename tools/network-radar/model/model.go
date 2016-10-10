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
	"github.com/cSploit/daemon/models"
	"github.com/op/go-logging"
	"golang.org/x/net/context"

	"net"
)

var log = logging.MustGetLogger("netowrk_radar.model")

func walkHostsIP(ctx context.Context, n *models.Network) <-chan net.IP {
	c := make(chan net.IP)

	go func() {
		for _, h := range n.GetHosts() {
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
	n := models.FindOrCreateNetwork(ipNet)

	return func(ctx context.Context) <-chan net.IP {
		return walkHostsIP(ctx, n)
	}
}
