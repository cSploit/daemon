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
	"fmt"
	ctxHelper "github.com/cSploit/daemon/helpers/ctx"
	"github.com/op/go-logging"
	"github.com/vektra/errors"
	"golang.org/x/net/context"
	"net"
)

var (
	log = logging.MustGetLogger("network-radar")
)

type NetworkRadar struct {
	Passive   bool
	Iface     *net.Interface
	Addresses []net.Addr
	ctx       context.Context
	Cancel    context.CancelFunc
	Receiver  HostReceiverFunc
	Fetcher   HostFetcher
}

func (nr *NetworkRadar) startProbing() error {
	var lastErr error
	var skipLoopback bool

	if nr.Fetcher == nil {
		return errors.New("Active scan requires an HostFetcher")
	}

	nr.Fetcher = nr.Fetcher.WithContext(nr.ctx)

	if len(nr.Addresses) == 0 {
		if nr.Iface != nil {
			nr.Addresses, lastErr = nr.Iface.Addrs()
		} else {
			nr.Addresses, lastErr = net.InterfaceAddrs()
		}

		if lastErr != nil {
			return lastErr
		}

		skipLoopback = true
	}

	lastErr = errors.New("no network to probe for")
	activated := 0

	for _, a := range nr.Addresses {

		ipNet, ok := a.(*net.IPNet)

		if !ok {
			log.Debugf("skipping non-ip address: <%T> %v", a, a)
			continue
		}

		if skipLoopback && ipNet.IP.IsLoopback() {
			continue
		}

		ctx := ctxHelper.WithIpNet(nr.ctx, ipNet)
		ctx = context.WithValue(ctx, fetcherKey, nr.Fetcher)

		if err := ProbeNetBIOS(ctx); err != nil {
			log.Error(err)
			lastErr = err
		} else {
			activated++
		}

		if err := ProbeKnownHosts(ctx); err != nil {
			log.Error(err)
			lastErr = err
		} else {
			activated++
		}
	}

	if activated == 0 {
		return fmt.Errorf("unable to start probers: %v", lastErr)
	}

	return nil
}

func (nr *NetworkRadar) Start() error {
	nr.ctx, nr.Cancel = context.WithCancel(context.Background())

	if nr.Iface != nil {
		nr.ctx = ctxHelper.WithIface(nr.ctx, *(nr.Iface))
	}

	a := analyzer{ctx: nr.ctx, Receiver: nr.Receiver}

	if err := a.Start(); err != nil {
		return err
	}

	if !nr.Passive {
		if err := nr.startProbing(); err != nil {
			nr.Cancel()
			return err
		}
	}

	return nil
}
