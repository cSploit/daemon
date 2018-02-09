package models

import (
	"github.com/cSploit/daemon/models/internal"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func init() {
	internal.RegisterModels(&ProcessJob{})
	registerJobKind(ProcessJobKind, "process")
}

const ProcessJobKind JobKind = 1

var commands = make(map[uint]*exec.Cmd)

//TODO: event system
var completed = make(map[uint]chan int)

type (
	ProcessJob struct {
		//TODO: hide job
		Job   Job
		JobId uint `gorm:"primary_key"`

		Command string `json:"command"`
		Args    string `json:"args"`

		//TODO: OutputHolder
		Output string `json:"output"`

		ExitStatus *int `json:"exit_status"`
	}

	ioManager struct {
		job   *ProcessJob
		stdin io.Writer
	}
)

func (m ioManager) Write(p []byte) (int, error) {
	m.job.Output += string(p)
	//TODO: save output asynchronously
	//TODO: save stdout and stderr separately but with correct order ( OutputHolder )
	if err := internal.Db.Model(m.job).Update("Output", m.job.Output).Error; err != nil {
		log.Error(err)
	}
	return len(p), nil
}

func (m ioManager) WriteToStdin(p []byte) (int, error) {
	return m.stdin.Write(p)
}

func (m ioManager) CloseStdin() (e error) {
	if closer, ok := m.stdin.(io.Closer); ok {
		e = closer.Close()
	}
	return
}

func (pj *ProcessJob) onStartFail(err error) {
	t := time.Now()
	status := 0
	db := internal.Db

	db.Model(pj).Updates(map[string]interface{}{
		"Output":     err.Error(),
		"ExitStatus": &status,
	})

	db.Model(&pj.Job).Update("FinishedAt", &t)

	pj.onDone()
}

func (pj *ProcessJob) onDone() {
	completed[pj.JobId] <- 0
}

func runCommand(pj ProcessJob, cmd *exec.Cmd) {
	statusCode := 0

	if err := cmd.Start(); err != nil {
		pj.onStartFail(err)
		return
	}

	err := cmd.Wait()
	end := time.Now()
	db := internal.Db

	log.Debugf("process %v exited: err=%v", pj, err)

	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			statusCode = status.ExitStatus()
			//TODO: signal and other goodies
		}
	} else if err != nil {
		log.Errorf("unexpected wait error %v", err)
	}

	if err := db.Model(&pj).Update("ExitStatus", &statusCode).Error; err != nil {
		log.Error(err)
	}

	if err := db.Model(&pj.Job).Update("FinishedAt", &end).Error; err != nil {
		log.Error(err)
	}

	pj.onDone()
}

func CreateProcessJob(command string, args ...string) (*ProcessJob, error) {

	name := command

	if len(args) > 0 {
		name += " " + strings.Join(args, " ")
	}

	j := Job{Name: name}

	pj := &ProcessJob{
		Command: command,
		Args:    strings.Join(args, string(0x17)),
		Job:     j,
	}

	if e := internal.Db.Create(pj).Error; e != nil {
		return nil, e
	}

	cmd := exec.Command(command, args...)

	iom := &ioManager{job: pj}

	if stdin, err := cmd.StdinPipe(); err != nil {
		log.Error(err)
		log.Warning("failed to attach process stdin")
		iom.stdin = ioutil.Discard
	} else {
		iom.stdin = stdin
	}

	cmd.Stdout = iom
	cmd.Stderr = iom

	completed[pj.JobId] = make(chan int, 1)
	commands[pj.JobId] = cmd

	go runCommand(*pj, cmd)

	return pj, nil
}

func FindProcessJob(id uint) (*ProcessJob, error) {
	j := &ProcessJob{}

	if e := internal.Db.Find(j, id).Error; e != nil {
		return nil, e
	}

	return j, nil
}

func (pj *ProcessJob) cmd() *exec.Cmd {
	return commands[pj.JobId]
}

func (pj *ProcessJob) ioManager() *ioManager {
	return pj.cmd().Stdout.(*ioManager)
}

// write to process stdin
func (pj *ProcessJob) Write(p []byte) (int, error) {
	return pj.ioManager().WriteToStdin(p)
}

// close process stdin
func (pj *ProcessJob) CloseInput() {
	pj.ioManager().CloseStdin()
}

// kill the job, runCommand will do the rest (I think)
func (pj *ProcessJob) Kill() (e error) {
	// Kill only if not completed
	if pj.ExitStatus == nil {
		// Retrieve cmd
		cmd := pj.cmd()

		// Kill the proc
		e = cmd.Process.Kill()

		if e != nil {
			log.Error(e)
		}
	}

	return
}
