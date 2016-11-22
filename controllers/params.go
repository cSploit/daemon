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
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const idLabel = "id"

func fetchId(c *gin.Context, entityName string, id *uint64) (e error) {
	if e = getId(c, entityName, id); e != nil {
		c.AbortWithError(http.StatusBadRequest, e)
	}

	return
}

func getId(c *gin.Context, entityName string, id *uint64) (e error) {
	var i uint64
	path := entityName + "_" + idLabel

	if s, exists := c.Params.Get(path); exists {
		if i, e = strconv.ParseUint(s, 10, 64); e == nil {
			*id = i
		}
	} else {
		e = errors.New(path + " not found")
	}
	return
}
