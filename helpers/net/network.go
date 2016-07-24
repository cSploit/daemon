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
package net

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/op/go-logging"
	"math"
	"net"
)

var (
	log = logging.MustGetLogger("helpers")
)

// gives the interface that is used to connect to an IP
func InterfaceForIp(ip net.IP) (net.Interface, error) {
	ifaces, err := net.Interfaces()

	if err != nil {
		return net.Interface{}, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()

		if err != nil {
			continue
		}

		for _, addr := range addrs {
			switch addr.(type) {
			case *net.IPNet:
				ipNet := addr.(*net.IPNet)
				if ipNet.Contains(ip) {
					return iface, nil
				}
			}
		}
	}

	return net.Interface{}, fmt.Errorf("Address %s unreachable", ip)
}

func GetInterfaceIP(iface net.Interface) (*net.IPNet, error) {
	return getIfaceIp(iface, false)
}

func GetInterfaceIPv4(iface net.Interface) (*net.IPNet, error) {
	return getIfaceIp(iface, true)
}

func GetMyEndpoints() ([]gopacket.Endpoint, error) {
	ifaces, err := net.Interfaces()

	if err != nil {
		return nil, err
	}

	var res = make([]gopacket.Endpoint, 0)

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()

		if err != nil {
			log.Error(err)
			continue
		}

		for _, addr := range addrs {
			switch addr.(type) {
			case *net.IPNet:
				ipNet := addr.(*net.IPNet)
				var et gopacket.EndpointType

				if ipNet.IP.To4() != nil {
					et = layers.EndpointIPv4
				} else {
					et = layers.EndpointIPv6
				}
				e := gopacket.NewEndpoint(et, ipNet.IP)
				res = append(res, e)
			default:
				log.Debugf("iface %s: got address <%T>: %v", iface, addr, addr)
			}
		}
	}

	return res, nil
}

func MacAddrToUInt(hw net.HardwareAddr) (uint64, error) {
	var val uint64
	var raw []byte

	if len(hw) < 8 {
		raw = make([]byte, 8)
		copy(raw[8-len(hw):], hw)
	} else {
		raw = hw
	}

	buf := bytes.NewReader(raw)
	err := binary.Read(buf, binary.BigEndian, &val)

	if err != nil {
		log.Warningf("unable to convert %v to uint64: %v", hw, err)
		return 0, err
	}

	return val, nil
}

func BuildBroadcastAddress(ipNet *net.IPNet) net.IP {
	res := ipNet.IP.Mask(ipNet.Mask)

	for i := 0; i < len(res); i++ {
		res[i] &= ipNet.Mask[i]
		res[i] |= ipNet.Mask[i] ^ 0xff
	}

	return res
}

// NextIP increase the passed ip by 1
func NextIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}

func CopyIP(ip net.IP) net.IP {
	// attempt to reduce size
	if ip4 := ip.To4(); ip4 != nil {
		ip = ip4
	}

	res := make(net.IP, len(ip))
	copy(res, ip)
	return res
}

func NumHosts(ipNet *net.IPNet) uint64 {
	ones, bits := ipNet.Mask.Size()
	zeros := float64(bits - ones)
	res := math.Pow(2, zeros) - 2
	res = math.Max(res, 0)

	return uint64(res)
}

// IPNetTo4 convert an IP Network to it's IPv4 short form.
// if the given IP network is not an IPv4 Network it returns nil
func IPNetTo4(ipNet *net.IPNet) *net.IPNet {
	if ip4 := ipNet.IP.To4(); ip4 != nil {
		return &net.IPNet{
			IP:   ip4,
			Mask: ipNet.Mask[len(ipNet.Mask)-4:],
		}
	}
	return nil
}

func getIfaceIp(iface net.Interface, ipv4Only bool) (*net.IPNet, error) {
	var addrs []net.Addr
	var err error

	if addrs, err = iface.Addrs(); err != nil {
		return nil, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			if net4 := IPNetTo4(ipnet); net4 != nil {
				ipnet = net4
			} else if ipv4Only {
				continue
			}
			return ipnet, nil
		}
	}

	return nil, fmt.Errorf("no IP addresses for interface '%s'", iface)
}
