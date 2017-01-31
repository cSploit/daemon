package network_radar

import (
	"net"
	"golang.org/x/net/context"
)

// fetch known hosts
type HostFetcher interface {
	WithContext(context.Context) HostFetcher
	WithNetwork(*net.IPNet) HostFetcher
	Find() <-chan net.IP
}