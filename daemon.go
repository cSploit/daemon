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
	"github.com/gin-gonic/gin"
	"gopkg.in/gin-contrib/cors.v1"
	"github.com/lair-framework/go-nmap"
	"github.com/op/go-logging"

	"flag"
	"github.com/cSploit/daemon/config"
	"github.com/ianschenck/envflag"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net"
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
	g, f := "google.com", "facebook.com"

	h1 := &models.Host{
		Name:   &g,
		IpAddr: "172.217.16.174",
	}
	h2 := &models.Host{
		Name:   &f,
		IpAddr: "31.13.76.68",
	}
	db := models.GetDbInstance()

	db.Create(h1)
	db.Create(h2)
}

func startRadars() {
	ifaces, err := net.Interfaces()

	if err != nil {
		panic(err)
	}

	for _, iface := range ifaces {

		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		i, err := models.FindIfaceByName(iface.Name)

		if err == gorm.ErrRecordNotFound {
			i, err = models.CreateIface(iface)
		}

		if err != nil {
			log.Error(err)
			continue
		}

		if job, err := controllers.IfaceScanz(i, &iface, config.Conf.Scan.Passive); err != nil {
			log.Error(err)
		} else {
			log.Infof("NetworkRadar succesfully started on interface %s: job#%d", iface.Name, job.ID)
		}
	}
}

func main() {
	flag.Parse()
	envflag.Parse()

	if err := config.Load(); err != nil {
		panic(err)
	}

	if err := models.Setup(); err != nil {
		log.Fatalf("unable to setup model: %v", err)
		panic("unable to setup model")
	}

	logging.SetBackend(logging.NewLogBackend(os.Stderr, "", 0))

	startRadars()

	router := gin.Default()

	router.Use(cors.Default()) //TODO: true CORS rules

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
		jobs := hc.NestedGroup(hosts, "/jobs")
		{
			controllers.JobController.Setup(jobs)
		}
	}
	networks := router.Group("networks")
	{
		nc := controllers.NetworkController
		nc.Setup(networks)

		jobs := nc.NestedGroup(networks, "/jobs")
		{
			controllers.JobController.Setup(jobs)
		}
	}
	ifaces := router.Group("ifaces")
	{
		ic := controllers.IfaceController
		ic.Setup(ifaces)
		actions := ic.NestedGroup(ifaces, "/")
		{
			actions.POST("scan", controllers.IfaceScan)
		}

		jobs := ic.NestedGroup(ifaces, "/jobs")
		{
			controllers.JobController.Setup(jobs)
		}
	}
	jobs := router.Group("jobs")
	{
		controllers.JobController.Setup(jobs)
	}

	router.Run(":8080")
}
