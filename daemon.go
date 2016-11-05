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
package main

import (
	"github.com/cSploit/daemon/controllers"
	"github.com/cSploit/daemon/models"
	"github.com/cSploit/daemon/tools/network-radar"
	"github.com/gin-gonic/gin"
	"github.com/lair-framework/go-nmap"
	"github.com/op/go-logging"

	"flag"
	"github.com/cSploit/daemon/config"
	"github.com/ianschenck/envflag"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"os"
)

var log = logging.MustGetLogger("daemon")

func loadScanFromFile(f string) error {
	xml, err := ioutil.ReadFile(f)

	if err != nil {
		return err
	}

	scan, err := nmap.Parse(xml)

	if err != nil {
		return err
	}

	db := models.GetDbInstance()

	for _, h := range scan.Hosts {
		db.Create(models.NewHost(h))
	}

	return nil
}

func initAllHostWithNetwork(ifName string, ipAddr string) {
	var hosts []models.Host

	db := models.GetDbInstance()
	network := models.NewNetwork("wlan0", "10.169.64.0/20")

	db.Find(&hosts)

	network.Hosts = hosts

	db.Create(network)
}

func addSomeRemoteHost() {
	h1 := &models.Host{
		Name:   null.StringFrom("google.com"),
		IpAddr: "172.217.16.174",
	}
	h2 := &models.Host{
		Name:   null.StringFrom("facebook.com"),
		IpAddr: "31.13.76.68",
	}
	db := models.GetDbInstance()

	db.Create(h1)
	db.Create(h2)
}

func main() {
	flag.Parse()
	envflag.Parse()

	if err := config.Load(); err != nil {
		panic(err)
	}

	var err = models.Setup()

	if err != nil {
		log.Fatalf("unable to setup model: %v", err)
		panic("unable to setup model")
	}

	logging.SetBackend(logging.NewLogBackend(os.Stderr, "", 0))

	err = loadScanFromFile("sample_nmap_out.xml")

	if err != nil {
		log.Errorf("Error loading nmap output: %v", err)
	}

	initAllHostWithNetwork("wlan0", "10.169.64.0/20")
	addSomeRemoteHost()

	router := gin.Default()

	hosts := router.Group("/hosts")
	{
		hc := controllers.HostsController
		hc.Setup(hosts)
		ports := hc.NestedGroup(hosts, "/ports")
		{
			controllers.PortsController.Setup(ports)
		}
		services := hc.NestedGroup(hosts, "/services")
		{
			controllers.ServicesController.Setup(services)
		}
	}
	networks := router.Group("networks")
	{
		controllers.NetworkController.Setup(networks)
	}

	nr := network_radar.NetworkRadar{}

	if err := nr.Start(); err != nil {
		log.Errorf("cannot start NetworkRadar: %v", err)
	}

	router.Run(":8080")
}
