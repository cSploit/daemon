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
	"github.com/lair-framework/go-nmap"
)

type Service struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Name    string `json:"name"`
	Product string `json:"product,omitempty"`
	Version string `json:"version,omitempty"`
	PortID  uint   `json:"-"`
}

func NewService(s nmap.Service) *Service {
	return &Service{Name: s.Name, Version: s.Version, Product: s.Product}
}

func (s *Service) FormatName() string {
	var res = s.Name

	if s.Product != "" {
		res += " - " + s.Product
	}

	if s.Version != "" {
		res += " ( " + s.Version + " )"
	}

	return res
}
