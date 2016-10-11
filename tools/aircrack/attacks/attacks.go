package attacks

import (
	"os"
	"time"
)

// To keep track of runned attacks
type Attack struct {
	Type    string `json:"type"`
	Target  string `json:"target"`
	Running bool   `json:"running"`
	Started string `json:"started at"`
	Stopped string `json:"stopped at"`
	process *os.Process
}

// Quick hack to edit process
func (a *Attack) Init(proc *os.Process) {
	a.process = proc
}

func (a *Attack) Stop() error {
	err := a.process.Kill()
	if err != nil {
		return err
	}

	a.Running = false
	a.Stopped = time.Now().String()

	return nil
}
