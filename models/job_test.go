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

	pj, err := CreateProcessJob("date")

	require.Nil(t, err)

	db.Create(&h)

	db.Model(&(pj.Job)).Association("Hosts").Append(h)

	require.Nil(t, db.Model(&h).Association("Jobs").Error)

	db.Model(&h).Association("Jobs").Find(&jobs)

	require.Equal(t, 1, len(jobs))
	require.Equal(t, pj.JobId, jobs[0].ID)
}

func TestRegisterJobKind(t *testing.T) {
	RegisterTestingT(t)

	v1, v2 := false, false
	k1 := JobKind(500)
	k2 := JobKind(500)

	f := func() {
		registerJobKind(k1, "good")
		v1 = true
		registerJobKind(k2, "fail")
		v2 = true
	}

	Expect(f).To(Panic())
	Expect(v1).To(BeTrue())
	Expect(v2).To(BeFalse())
}
