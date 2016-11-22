package views

import "github.com/cSploit/daemon/models"

type clientShowView struct {
	models.Client
	Iface interface{} `json:"iface"`
	Jobs  interface{} `json:"jobs"`
}

func ClientIndex(arg interface{}) interface{} {
	return arg
}

func ClientShow(arg interface{}) interface{} {
	client := arg.(models.Client)

	ifaces := []models.Iface{client.Iface}

	return &clientShowView{
		Client: client,
		Iface:  IfaceIndex(ifaces),
		Jobs:   JobIndex(client.Jobs),
	}
}
