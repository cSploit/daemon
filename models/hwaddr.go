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
	"errors"
	netHelper "github.com/cSploit/daemon/helpers/net"
	"github.com/lair-framework/go-nmap"
	"net"
)

type HwAddr struct {
	ID       uint
	RawValue uint64 `gorm:"index"`
	Addr     string
	Vendor   string
	HostId   uint
	Host     *Host `json:"-"`
}

func NewHwAddr(a interface{}) (*HwAddr, error) {
	var vendor string
	var str string
	var id uint64
	var err error

	switch a.(type) {
	default:
		log.Errorf("unexpected type %T", a)
		return &HwAddr{}, errors.New("unexpected type")
	case nmap.Address:
		str = a.(nmap.Address).Addr
		vendor = a.(nmap.Address).Vendor
		id, err = MACStringToRaw(str)
	case string:
		str = a.(string)
		id, err = MACStringToRaw(str)
	case net.HardwareAddr:
		str = a.(net.HardwareAddr).String()
		id, err = netHelper.MacAddrToUInt(a.(net.HardwareAddr))
	case *net.HardwareAddr:
		str = a.(*net.HardwareAddr).String()
		id, err = netHelper.MacAddrToUInt(*(a.(*net.HardwareAddr)))
	}

	if err != nil {
		log.Error("bad mac address: ", err)
		return &HwAddr{}, err
	}

	return &HwAddr{RawValue: id, Addr: str, Vendor: vendor}, nil
}

func MACStringToRaw(str string) (uint64, error) {
	hw, err := net.ParseMAC(str)
	if err != nil {
		log.Warning("Bad MAC Address: ", err)
		return 0, err
	}
	return netHelper.MacAddrToUInt(hw)
}
