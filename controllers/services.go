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

var ServicesController = Controller{
	EntityName: "service",
	Index:      servicesIndex,
	Show:       servicesShow,
}

func servicesIndex(c *gin.Context) {
	var host_id uint64
	var services []models.Service

	db := models.GetDbInstance()

	if fetchId(c, "host", &host_id) != nil {
		return
	}

	dbRes := db.Joins("JOIN ports ON port_id = ports.id").Where("host_id = ?", host_id).Find(&services)

	renderView(c, views.ServiceIndex, services, dbRes)
}

func servicesShow(c *gin.Context) {
	var host_id uint64
	var id uint64
	var svc models.Service

	db := models.GetDbInstance()

	if fetchId(c, "host", &host_id) != nil {
		return
	}

	if fetchId(c, "service", &id) != nil {
		return
	}

	dbRes := db.Joins("JOIN ports ON port_id = ports.id").Where("host_id = ?", host_id).Find(&svc, "services.id = ?", id)

	renderView(c, views.ServiceShow, svc, dbRes)
}
