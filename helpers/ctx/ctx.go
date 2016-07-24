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
package ctx

import (
	"golang.org/x/net/context"
	"net"
)

type key int

const ipNetKey key = 0
const ifaceKey key = 1

func WithIpNet(ctx context.Context, ipNet *net.IPNet) context.Context {
	return context.WithValue(ctx, ipNetKey, ipNet)
}

func GetIpNet(ctx context.Context) *net.IPNet {
	return ctx.Value(ipNetKey).(*net.IPNet)
}

func WithIface(ctx context.Context, iface net.Interface) context.Context {
	return context.WithValue(ctx, ifaceKey, iface)
}

func GetIface(ctx context.Context) net.Interface {
	return ctx.Value(ifaceKey).(net.Interface)
}

func HaveIface(ctx context.Context) bool {
	return ctx.Value(ifaceKey) != nil
}
