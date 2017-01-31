package network_radar

import "net"

// Receive notification about seen hosts
// hwAddr is the link address of the host, ipAddr it's IP one and name it's DNS or NetBIOS name if available
type HostReceiverFunc func(hwAddr net.HardwareAddr, ipAddr net.IP, name *string)


