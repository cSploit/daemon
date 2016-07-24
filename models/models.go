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
	"github.com/cSploit/daemon/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func Setup() error {
	var err error

	var models = []interface{}{&Host{}, &Port{}, &Service{},
		&HwAddr{}, &Network{}}

	db, err = gorm.Open(config.Conf.Db.Dialect, config.Conf.Db.Args...)

	if err != nil {
		return err
	}

	db = db.Debug().
		DropTableIfExists(models...).
		AutoMigrate(models...)

	return nil
}

func GetDbInstance() gorm.DB {
	return *db
}
