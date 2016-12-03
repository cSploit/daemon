package models

import (
	"github.com/cSploit/daemon/models/internal"
	"github.com/ianschenck/envflag"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJobHosts(t *testing.T) {
	envflag.Parse()
	internal.OpenDbForTests()

	h := Host{IpAddr: "test"}
	db := internal.Db
	var jobs []Job

	pj, _ := CreateProcessJob("date")
	db.Create(&h)

	db.Model(&(pj.Job)).Association("Hosts").Append(h)

	require.Nil(t, db.Model(&h).Association("Jobs").Error)

	db.Model(&h).Association("Jobs").Find(&jobs)

	require.Equal(t, 1, len(jobs))
	require.Equal(t, pj.JobId, jobs[0].ID)
}
