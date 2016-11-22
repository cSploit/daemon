package models

import (
	"github.com/cSploit/daemon/models/internal"
)

func init() {
	internal.RegisterModels(&RadarJob{})
	registerJobKind(RadarJobKind, "radar")
}

//TODO atExit(markAllRadarsAsFinished)

const RadarJobKind JobKind = 2

type RadarJob struct {
	internal.Base
	Job   Job
	JobId uint
}
