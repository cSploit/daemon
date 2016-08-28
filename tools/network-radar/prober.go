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
	netHelper "github.com/cSploit/daemon/helpers/net"
	"github.com/cSploit/daemon/tools/network-radar/model"
	"github.com/cSploit/daemon/tools/network-radar/netbios"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/vektra/errors"
	"golang.org/x/net/context"
	"math"
	"net"
	"sync"
	"time"
)

// hosts4sock is preferred amount of hosts t probe per open socket
const hosts4sock = 32

// maxSocks is maximum number of opened sockets
const maxSocks = 32

func walkNetwork(ctx context.Context) <-chan net.IP {
	c := make(chan net.IP)
	ipNet := ctxHelper.GetIpNet(ctx)

	go func() {
		defer close(c)

		ip := ipNet.IP.Mask(ipNet.Mask)
		bcast := netHelper.BuildBroadcastAddress(ipNet)

		// single address network
		if ip.Equal(bcast) {
			c <- ip
			return
		}

		for netHelper.NextIP(ip); !ip.Equal(bcast); netHelper.NextIP(ip) {
			res := netHelper.CopyIP(ip)
			select {
			case c <- res:
			case <-ctx.Done():
				return
			}
		}
	}()

	return c
}

func loopKnownHosts(ctx context.Context, loopDuration time.Duration, walker model.KnownHostsIPWalker) <-chan net.IP {
	c := make(chan net.IP)
	ticker := time.NewTicker(loopDuration)

	pipe := func(in <-chan net.IP) {
		for ip := range in {
			select {
			case c <- ip:
			case <-ctx.Done():
				return
			}
		}
	}

	go func() {
		var warned bool

		defer close(c)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				start := time.Now()
				pipe(walker(ctx))
				elapsed := time.Since(start)
				if !warned && elapsed > loopDuration {
					warned = true
					log.Warningf("Want to walk the known hosts every %v but we took %v",
						loopDuration, elapsed)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return c
}

func nbProbe(ctx context.Context, c <-chan net.IP) (<-chan struct{}, error) {
	laddr := &net.UDPAddr{IP: net.IPv4zero}
	nbConn, err := net.ListenUDP("udp", laddr)

	if err != nil {
		return nil, err
	}

	done := make(chan struct{})

	go func() {
		defer nbConn.Close()
		defer close(done)
		for {
			select {
			case ip, more := <-c:
				if !more {
					return
				}
				err := netbios.SendQuery(nbConn, ip)
				if err != nil {
					log.Error(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return done, nil
}

func ipv4ArpRequestGenerator(ctx context.Context, c <-chan net.IP) <-chan gopacket.SerializeBuffer {
	iface := ctxHelper.GetIface(ctx)
	ipNet := ctxHelper.GetIpNet(ctx)

	srcIp := ipNet.IP.To4()

	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	eth := &layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: []byte(srcIp),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}

	out := make(chan gopacket.SerializeBuffer)

	go func() {
		defer close(out)

		for {
			select {
			case ip, more := <-c:
				if !more {
					return
				}

				buf := gopacket.NewSerializeBuffer()

				arp.SourceProtAddress = []byte(ip.To4())
				gopacket.SerializeLayers(buf, opts, eth, arp)

				log.Debugf("ARP Request length: %d", len(buf.Bytes()))

				out <- buf
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func interfaceWriter(ctx context.Context, c <-chan gopacket.SerializeBuffer) error {
	iface := ctxHelper.GetIface(ctx)
	handle, err := pcap.OpenLive(iface.Name, 0, true, pcap.BlockForever)

	if err != nil {
		return err
	}

	go func() {
		defer handle.Close()

		for {
			select {
			case buf, more := <-c:
				if !more {
					return
				}
				handle.WritePacketData(buf.Bytes())
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func mergeBufs(ctx context.Context, chans ...<-chan gopacket.SerializeBuffer) <-chan gopacket.SerializeBuffer {
	var wg sync.WaitGroup

	out := make(chan gopacket.SerializeBuffer)

	pipe := func(c <-chan gopacket.SerializeBuffer) {
		defer wg.Done()

		for buf := range c {
			select {
			case out <- buf:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(len(chans))

	for _, c := range chans {
		go pipe(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func tryToReduceSize(ipNet *net.IPNet) *net.IPNet {
	if net4 := netHelper.IPNetTo4(ipNet); net4 != nil {
		return net4
	}
	return ipNet
}

func ProbeNetBIOS(ctx context.Context) error {
	var lastErr error

	ipNet := ctxHelper.GetIpNet(ctx)

	N := netHelper.NumHosts(ipNet)

	if N == 0 {
		return fmt.Errorf("Network '%s' is empty", ipNet)
	}

	ctx, cancel := context.WithCancel(ctx)
	NSenders := int(math.Ceil(float64(N) / hosts4sock))
	NSenders = int(math.Min(float64(NSenders), maxSocks))
	ips := walkNetwork(ctx)
	activated := 0

	log.Infof("starting NetBIOS prober for network '%s' {N: %d, NSenders: %d }", ipNet, N, NSenders)

	for i := 0; i < NSenders; i++ {
		_, lastErr = nbProbe(ctx, ips)

		if lastErr != nil {
			log.Error(lastErr)
			continue
		}

		activated++
	}

	if activated == 0 {
		cancel()
		return fmt.Errorf("Cannot create probes: %v", lastErr)
	}

	return nil
}

func ProbeKnownHosts(ctx context.Context) error {
	ipNet := ctxHelper.GetIpNet(ctx)
	ctx, cancel := context.WithCancel(ctx)

	if net4 := netHelper.IPNetTo4(ipNet); net4 == nil {
		return errors.New("IPv6 not implemented yet")
	} else {
		ipNet = net4
		ctx = ctxHelper.WithIpNet(ctx, ipNet)
	}

	N := netHelper.NumHosts(ipNet)

	if N == 0 {
		return fmt.Errorf("Network '%s' is empty", ipNet)
	}

	iface, err := netHelper.InterfaceForIp(ipNet.IP)

	if err != nil {
		return err
	}

	ctx = ctxHelper.WithIface(ctx, iface)

	walker := model.NewKnownHostsWalker(ipNet)
	ips := loopKnownHosts(ctx, time.Second, walker)

	NGens := int(math.Ceil(float64(N) / hosts4sock))
	NGens = int(math.Min(float64(NGens), maxSocks))
	var buffChannels []<-chan gopacket.SerializeBuffer

	log.Infof("starting ARP prober for network '%s' {N: %d, NGens: %d, iface : %v}", ipNet, N, NGens, iface)

	for i := 0; i < NGens; i++ {
		bufc := ipv4ArpRequestGenerator(ctx, ips)

		//TODO: ipv6NeighborRequestGenerator

		buffChannels = append(buffChannels, bufc)
	}

	bufs := mergeBufs(ctx, buffChannels...)

	if err := interfaceWriter(ctx, bufs); err != nil {
		cancel()
		return err
	}

	return nil
}
