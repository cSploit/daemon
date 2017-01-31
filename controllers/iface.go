package controllers

import (
	"github.com/cSploit/daemon/helpers"
	"github.com/cSploit/daemon/models"
	"github.com/cSploit/daemon/tools/network-radar"
	"github.com/cSploit/daemon/views"
	"github.com/gin-gonic/gin"
	"github.com/vektra/errors"
	"net"
	"net/http"
	"strconv"
)

var IfaceController = Controller{
	EntityName: "iface",
	Index:      ifaceIndex,
	Show:       ifaceShow,
	//TODO: methods ( POST /ifaces/1/scan ) [ gutron ]
}

var radars = make(map[uint]network_radar.NetworkRadar)

func ifaceIndex(c *gin.Context) {
	var ifaces []models.Iface

	db := models.GetDbInstance().Find(&ifaces)

	renderView(c, views.IfaceIndex, ifaces, db)
}

func ifaceShow(c *gin.Context) {
	var id uint64

	if fetchId(c, "iface", &id) != nil {
		return
	}

	iface := models.Iface{}

	db := models.GetDbInstance().Find(&iface, id)

	renderView(c, views.IfaceShow, iface, db)
}

func IfaceScan(c *gin.Context) {
	var id uint64
	var err error
	var job *models.Job
	var passive = false
	var netIface *net.Interface

	if fetchId(c, "iface", &id) != nil {
		return
	}

	iface := &models.Iface{}

	db := models.GetDbInstance().
		Preload("Jobs", "? = type AND finished_at = NULL", models.RadarJobKind).
		Find(iface, id)

	if db.Error != nil {
		goto done
	}

	if len(iface.Jobs) > 0 {
		db.Error = errors.Format("Radar already running")
		goto done
	}

	if arg, haveIt := c.GetPostForm("passive"); haveIt {
		if passive, err = strconv.ParseBool(arg); err != nil {
			db.Error = err
			goto error
		}
	}

	if netIface, err = net.InterfaceByName(iface.Name); err != nil {
		goto error
	}

	if job, err = IfaceScanz(iface, netIface, passive); err == nil {
		goto done
	}

error:
	log.Error(err)
	c.AbortWithStatus(http.StatusInternalServerError)

	return

done:

	renderView(c, views.JobShow, job, db)
}

// just an hack, will improve it when will switch to gutron
func IfaceScanz(model *models.Iface, iface *net.Interface, passive bool) (*models.Job, error) {

	nr := network_radar.NetworkRadar{
		Iface:    iface,
		Passive:  passive,
		Receiver: models.NotifyHostSeen,
		Fetcher:  helpers.BaseFetcher,
	}

	if err := nr.Start(); err != nil {
		return nil, err
	}

	job := &models.Job{}
	radarJob := &models.RadarJob{Job: *job}

	job.Ifaces = append(job.Ifaces, *model)

	if err := models.GetDbInstance().Save(radarJob).Error; err != nil {
		return nil, err
	}

	radars[radarJob.ID] = nr

	job.Radar = radarJob //FIXME: is this necessary ?

	return job, nil
}
