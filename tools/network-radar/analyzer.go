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
	"net"
	"sync"

	ctxHelper "github.com/cSploit/daemon/helpers/ctx"
	netHelper "github.com/cSploit/daemon/helpers/net"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/net/context"
)

var (
	localEndpoints = struct {
		sync.RWMutex
		Points []gopacket.Endpoint
	}{}
)

type analyzer struct {
	ctx context.Context
	// networks that the analyzer will deal with
	WatchedNetworks []*net.IPNet
	Receiver        HostReceiverFunc
}

func init() {
	//TODO spawn endpoints poller
	res, err := netHelper.GetMyEndpoints()

	if err != nil {
		log.Error(err)
		return
	}

	localEndpoints.Lock()
	localEndpoints.Points = res
	localEndpoints.Unlock()
}

func (a *analyzer) isInternal(ip net.IP) bool {
	for _, ipNet := range a.WatchedNetworks {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

func (a *analyzer) onPacket(pkt gopacket.Packet) {
	//TODO: Application Layer ( NetBIOS )
	if pkt.NetworkLayer() != nil {
		a.analyzeNetworkPkt(pkt)
	} else if pkt.LinkLayer() != nil {
		a.analyzeLinkPkt(pkt)
	}
}

func isOurEndpoint(e gopacket.Endpoint) bool {
	localEndpoints.RLock()
	defer localEndpoints.RUnlock()

	for _, ee := range localEndpoints.Points {
		if ee == e {
			return true
		}
	}
	return false
}

func (a *analyzer) analyzeLinkPkt(pkt gopacket.Packet) {
	if arpLayer := pkt.Layer(layers.LayerTypeARP); arpLayer != nil {
		a.analyzeARP(pkt)
	}
}

func (a *analyzer) analyzeARP(pkt gopacket.Packet) {
	ll := pkt.LinkLayer()
	flow := ll.LinkFlow()

	if isOurEndpoint(flow.Src()) {
		log.Debugf("skipping sent ARP packet")
		return
	}

	if a.Receiver == nil {
		log.Debugf("Receiver is null, ARP packet lost")
		return
	}

	arp := pkt.Layer(layers.LayerTypeARP).(*layers.ARP)

	hwAddr := net.HardwareAddr(flow.Src().Raw())
	ipAddr := net.IP(arp.SourceProtAddress)

	go a.Receiver(hwAddr, ipAddr, nil)
}

func (a *analyzer) analyzeNetworkPkt(pkt gopacket.Packet) {
	ll := pkt.LinkLayer()
	nl := pkt.NetworkLayer()

	llSrc, llDst := ll.LinkFlow().Endpoints()
	nlSrc, nlDst := nl.NetworkFlow().Endpoints()

	var lle, nle gopacket.Endpoint

	if !isOurEndpoint(nlSrc) {
		lle = llSrc
		nle = nlSrc
	} else {
		lle = llDst
		nle = nlDst
	}

	if nle.EndpointType() != layers.EndpointIPv4 && nle.EndpointType() != layers.EndpointIPv6 {
		log.Debugf("skipping non-ip packet: %v", pkt)
		return
	}

	hwAddr := net.HardwareAddr(lle.Raw())
	ipAddr := net.IP(nlSrc.Raw())

	if a.isInternal(ipAddr) || netHelper.IsPrivate(ipAddr) {
		go a.Receiver(hwAddr, ipAddr, nil)
	}
}

// start sniffing and analyzing packets
func (a *analyzer) Start() error {
	ifName := "any"

	if ctxHelper.HaveIface(a.ctx) {
		ifName = ctxHelper.GetIface(a.ctx).Name
	}

	if len(a.WatchedNetworks) == 0 {
		networks, err := netHelper.GetAttachedIpNetworks()

		if err != nil {
			return err
		}

		a.WatchedNetworks = networks
	}

	handle, err := pcap.OpenLive(ifName, 1024, true, pcap.BlockForever)

	if err != nil {
		return err
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())

	go func() {
		defer handle.Close()
		for {
			select {
			case p, more := <-source.Packets():
				if !more {
					return
				}
				//TODO: use workers
				a.onPacket(p)
			case <-a.ctx.Done():
				return
			}
		}
	}()

	return nil
}
