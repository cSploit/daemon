package models

import (
	"github.com/cSploit/daemon/models/internal"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	internal.OpenDbForTests()
}

func TestJobHosts(t *testing.T) {
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
