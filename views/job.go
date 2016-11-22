package views

import "github.com/cSploit/daemon/models"

type jobShowView struct {
	models.Job
	Aps      interface{} `json:"aps"`
	Clients  interface{} `json:"clients"`
	Hosts    interface{} `json:"hosts"`
	Networks interface{} `json:"networks"`
	Ifaces   interface{} `json:"ifaces"`
	Process  interface{} `json:"process"`
}

type processJobShowView struct {
	models.ProcessJob
	hiddenJob string `json:"job,omitempty"`
}

func JobIndex(arg interface{}) interface{} {
	return arg
}

func JobShow(arg interface{}) interface{} {
	job := arg.(models.Job)

	return &jobShowView{
		Job:      job,
		Aps:      ApIndex(job.Aps),
		Clients:  ClientIndex(job.Clients),
		Hosts:    HostsIndex(job.Hosts),
		Networks: NetworkIndex(job.Networks),
		Ifaces:   IfaceIndex(job.Ifaces),
		Process:  processJobShow(job.Process),
	}
}

func processJobShow(pj *models.ProcessJob) interface{} {
	if pj == nil {
		return nil
	}

	return &processJobShowView{ProcessJob: *pj}
}
