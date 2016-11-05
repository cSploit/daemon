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
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
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
var logger *logging.Logger = logging.MustGetLogger("config")

func init() {
	fpath := "./config.json"

	if envPath := os.Getenv("CSPLOIT_CONFIG"); envPath != "" {
		fpath = envPath
	}

	config_file, err := ioutil.ReadFile(fpath)
	if err != nil {
		logger.Fatal("config file not found", err)
	}

	err = json.Unmarshal(config_file, &Conf)

	if err != nil {
		logger.Fatal("config format error", err)
	}
}
