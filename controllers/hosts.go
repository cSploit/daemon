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

var HostsController = Controller{
	EntityName: "host",
	Index:      hostsIndex,
	Show:       hostsShow,
}

func hostsIndex(c *gin.Context) {

	hosts := make([]models.Host, 0)
	db := models.GetDbInstance()

	dbRes := db.Preload("Ports", "state = ?", "open").Find(&hosts)

	renderView(c, views.HostsIndex, hosts, dbRes)
}

func hostsShow(c *gin.Context) {
	var id uint64

	if err := fetchId(c, "host", &id); err != nil {
		return
	}

	db := models.GetDbInstance()

	var host models.Host

	dbRes := db.Preload("Ports").Preload("Ports.Service").
		Preload("Network").Find(&host, id)

	renderView(c, views.HostsShow, host, dbRes)
}
