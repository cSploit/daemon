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
package controllers

import (
	"github.com/cSploit/daemon/models"
	"github.com/cSploit/daemon/views"
	"github.com/gin-gonic/gin"
)

var PortsController = Controller{
	EntityName: "port",
	Index:      portsIndex,
	Show:       portsShow,
}

func portsIndex(c *gin.Context) {
	var ports []models.Port
	var host_id uint64

	db := models.GetDbInstance()

	if fetchId(c, "host", &host_id) != nil {
		return
	}

	dbRes := db.Preload("Service").Where("host_id = ?", host_id).Find(&ports)

	renderView(c, views.PortIndex, ports, dbRes)
}

func portsShow(c *gin.Context) {
	var port models.Port
	var host_id uint64
	var id uint64

	db := models.GetDbInstance()

	if fetchId(c, "host", &host_id) != nil {
		return
	}

	if fetchId(c, "port", &id) != nil {
		return
	}

	dbRes := db.Preload("Service").Where("host_id = ?", host_id).Find(&port, id)

	renderView(c, views.PortShow, port, dbRes)
}
