package views

import "github.com/cSploit/daemon/models"

type ifaceShowView struct {
	models.Iface
	Aps     interface{} `json:"aps"`
	Clients interface{} `json:"clients"`
	Jobs    interface{} `json:"jobs"`
}

func IfaceIndex(arg interface{}) interface{} {
	return arg
}

func IfaceShow(arg interface{}) interface{} {
	iface := arg.(models.Iface)

	return &ifaceShowView{
		Iface:   iface,
		Aps:     ApIndex(iface.Aps),
		Clients: ClientIndex(iface.Clients),
		Jobs:    JobIndex(iface.Jobs),
	}
}
