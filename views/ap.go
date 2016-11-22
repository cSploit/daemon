package views

import "github.com/cSploit/daemon/models"

type apShowView struct {
	models.AP
	IfaceView interface{} `json:"iface"`
	JobsView  interface{} `json:"jobs"`
}

func ApIndex(arg interface{}) interface{} {
	return arg
}

func ApShow(arg interface{}) interface{} {
	ap := arg.(models.AP)

	view := apShowView{
		AP:        ap,
		IfaceView: IfaceIndex(ap.Iface),
		JobsView:  JobIndex(ap.Jobs),
	}

	return view
}
