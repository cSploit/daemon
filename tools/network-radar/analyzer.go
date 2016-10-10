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
	"github.com/cSploit/daemon/models"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/net/context"
)

var (
	myEndpoints     []gopacket.Endpoint
	myEndpointsLock sync.RWMutex
)

func init() {
	//TODO spawn endpoints poller
	res, err := netHelper.GetMyEndpoints()

	if err != nil {
		log.Error(err)
		return
	}

	myEndpointsLock.Lock()
	myEndpoints = res
	myEndpointsLock.Unlock()
}

func onPacket(pkt gopacket.Packet) {

	log.Debugf("onPacket(%v)", pkt)

	if pkt.NetworkLayer() != nil {
		log.Debugf("received network packet: %v", pkt)
		analyzeNetworkPkt(pkt)
	} else if pkt.LinkLayer() != nil {
		analyzeLinkPkt(pkt)
	}
}

func isOurEndpoint(e gopacket.Endpoint) bool {
	myEndpointsLock.RLock()
	defer myEndpointsLock.RUnlock()

	for _, ee := range myEndpoints {
		if ee == e {
			return true
		}
	}
	return false
}

func analyzeLinkPkt(pkt gopacket.Packet) {
	if arpLayer := pkt.Layer(layers.LayerTypeARP); arpLayer != nil {
		analyzeARP(pkt)
	}
}

func analyzeARP(pkt gopacket.Packet) {
	ll := pkt.LinkLayer()
	flow := ll.LinkFlow()

	if isOurEndpoint(flow.Src()) {
		log.Debugf("skipping sent ARP packet")
		return
	}

	arp := pkt.Layer(layers.LayerTypeARP).(*layers.ARP)

	log.Debugf("received an ARP packet: %v", arp)
}

func analyzeNetworkPkt(pkt gopacket.Packet) {
	ll := pkt.LinkLayer()
	nl := pkt.NetworkLayer()

	llSrc, llDst := ll.LinkFlow().Endpoints()
	nlSrc, nlDst := nl.NetworkFlow().Endpoints()

	var lle, nle gopacket.Endpoint

	if !isOurEndpoint(llSrc) {
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

	if err := models.NotifyHostSeen(hwAddr, ipAddr, ""); err != nil {
		log.Error(err)
	}
}

// start sniffing and analyzing packets
func startAnalyze(ctx context.Context) error {
	ifName := "any"

	if ctxHelper.HaveIface(ctx) {
		ifName = ctxHelper.GetIface(ctx).Name
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
				onPacket(p)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
