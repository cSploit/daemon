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
	"github.com/cSploit/daemon/views"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func renderView(c *gin.Context, render views.RenderFunc, args interface{}, dbResult *gorm.DB) {

	var err error
	var res interface{}

	if dbResult != nil {
		if dbResult.RecordNotFound() {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else if dbResult.Error != nil {
			err = dbResult.Error
			goto error
		}
	}

	res = render(args)

	c.JSON(http.StatusOK, res)

	return

error:
	log.Error(err)
	c.AbortWithError(http.StatusInternalServerError, err)
}
