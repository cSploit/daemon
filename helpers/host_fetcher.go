package helpers

import (
	"golang.org/x/net/context"
	"net"
	"github.com/cSploit/daemon/models"
	nr "github.com/cSploit/daemon/tools/network-radar"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("helpers")

type HostFetcher struct {
	network *models.Network
	ctx     context.Context
}

func (hf HostFetcher) WithNetwork(ipNet *net.IPNet) nr.HostFetcher {
	i := models.FindOrCreateNetwork(ipNet)
	return HostFetcher{network: i, ctx: hf.ctx}
}

func (hf HostFetcher) WithContext(ctx context.Context) nr.HostFetcher {
	return HostFetcher{network: hf.network, ctx: ctx}
}

func (hf HostFetcher) Find() <-chan net.IP {
	c := make(chan net.IP)
	var hosts []models.Host

	if hf.network != nil {
		hosts = hf.network.Hosts
	} else if err := models.GetDbInstance().Find(&hosts).Error; err != nil {
		log.Error(err)
		return c
	}

	go func() {
		for _, h := range hosts {
			ip := net.ParseIP(h.IpAddr)

			if ip == nil {
				log.Warningf("unable to parse ip '%s' for host %s", h.IpAddr, h)
				continue
			}

			select {
			case c <- ip:
			case <-hf.ctx.Done():
				return
			}
		}
	}()

	return c
}

var BaseFetcher = HostFetcher{}
