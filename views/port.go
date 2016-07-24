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
package views

import "github.com/cSploit/daemon/models"

type portIdxElem struct {
	models.Port
	HideHostId  string `json:"host_id,omitempty"`
	ServiceID   uint   `json:"service_id,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
}

type portShowView struct {
	models.Port
	Service interface{} `json:"service,omitempty"`
}

func PortIndex(args interface{}) interface{} {
	ports := args.([]models.Port)
	res := make([]portIdxElem, len(ports))

	for i, p := range ports {
		var svc string
		var svc_id uint

		if p.Service != nil {
			svc = p.Service.FormatName()
			svc_id = p.Service.ID
		}

		res[i] = portIdxElem{
			Port:        p,
			ServiceName: svc,
			ServiceID:   svc_id,
		}
	}

	return res
}

func PortShow(arg interface{}) interface{} {
	port := arg.(models.Port)

	view := &portShowView{Port: port}

	if port.Service != nil {
		view.Service = ServiceShow(*port.Service)
	}

	return view
}
