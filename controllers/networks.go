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

var NetworkController = Controller{
	EntityName: "network",
	Index:      networksIndex,
	Show:       networksShow,
}

func networksIndex(c *gin.Context) {
	db := models.GetDbInstance()
	var networks []models.Network

	dbRes := db.Find(&networks)

	renderView(c, views.NetworkIndex, networks, dbRes)
}

func networksShow(c *gin.Context) {
	db := models.GetDbInstance()
	var network models.Network
	var id uint64

	if fetchId(c, "network", &id) != nil {
		return
	}

	dbRes := db.Preload("Hosts").Find(&network, id)

	renderView(c, views.NetworkShow, network, dbRes)
}
