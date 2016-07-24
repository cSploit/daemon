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

/*
usually you'll get some error due to the high rate of emitting packets

errno: operation not permitted
dmesg: nf_conntrack: table full, dropping packet
solution: increase /proc/sys/net/ipv4/netfilter/ip_conntrack_max

errno: invalid argument
dmesg: neighbour: arp_cache: neighbor table overflow!
solution: http://www.cyberciti.biz/faq/centos-redhat-debian-linux-neighbor-table-overflow/
*/

import (
	ctxHelper "github.com/cSploit/daemon/helpers/ctx"
	"golang.org/x/net/context"
	"net"
	"sync"
	"testing"
)

func BenchmarkNbProber24_1(b *testing.B) {
	_, ipNet, _ := net.ParseCIDR("127.0.0.1/24")

	benchOne(b, ipNet, 1)
}

func BenchmarkNbProber24_4(b *testing.B) {
	_, ipNet, _ := net.ParseCIDR("127.0.0.1/24")

	benchOne(b, ipNet, 4)
}

func BenchmarkNbProber24_8(b *testing.B) {
	_, ipNet, _ := net.ParseCIDR("127.0.0.1/24")

	benchOne(b, ipNet, 8)
}

func benchOne(b *testing.B, ipNet *net.IPNet, NSenders int) {
	ctx := context.Background()
	ctx = ctxHelper.WithIpNet(ctx, ipNet)
	wg := sync.WaitGroup{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ips := walkNetwork(ctx)
		for j := 0; j < NSenders; j++ {
			done, err := nbProbe(ctx, ips)
			if err != nil {
				panic(err)
			}
			wg.Add(1)
			go func() {
				<-done
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
