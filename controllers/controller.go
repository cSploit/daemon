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

import "github.com/gin-gonic/gin"

type Controller struct {
	EntityName                          string
	Create, Index, Show, Update, Delete gin.HandlerFunc
}

func (ctrl Controller) Setup(r gin.IRouter) {
	id_path := "/:" + ctrl.EntityName + "_" + idLabel

	if ctrl.Index != nil {
		r.GET("/", ctrl.Index)
	}
	if ctrl.Create != nil {
		r.POST("/", ctrl.Create)
	}
	if ctrl.Show != nil {
		r.GET(id_path, ctrl.Show)
	}
	if ctrl.Update != nil {
		r.PATCH(id_path, ctrl.Update)
		r.PUT(id_path, ctrl.Update)
	}
	if ctrl.Delete != nil {
		r.DELETE(id_path, ctrl.Delete)
	}
}

func (ctrl Controller) NestedGroup(r gin.IRouter, relativePath string, args ...gin.HandlerFunc) *gin.RouterGroup {
	id_path := "/:" + ctrl.EntityName + "_" + idLabel
	return r.Group(id_path+relativePath, args...)
}
