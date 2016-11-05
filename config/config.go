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
package config

import (
	"encoding/json"
	"io/ioutil"
)

type DbConfig struct {
	Dialect string        `json:"dialect"`
	Args    []interface{} `josn:"args"`
}

type Config struct {
	Db DbConfig `json:"db"`
}

// global configuration object
var Conf Config

func Load(fpath string) error {
	var content []byte
	var err error

	if content, err = ioutil.ReadFile(fpath); err != nil {
		return err
	}

	if err = json.Unmarshal(content, &Conf); err != nil {
		return err
	}

	return nil
}
