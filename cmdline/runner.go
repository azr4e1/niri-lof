package cmdline

import (
	"fmt"
	"os/exec"
)

type Runner struct{}

func (r Runner) GetJSON() ([]byte, error) {
	niriMsg := exec.Command("niri", "msg", "-j", "windows")

	return niriMsg.Output()
}

func (r Runner) Focus(winID int) error {
	niriMsg := exec.Command("niri", "msg", "action", "focus-window", "--id", fmt.Sprintf("%d", winID))

	return niriMsg.Run()
}

func (r Runner) Spawn(cmd string) error {
	niriMsg := exec.Command("niri", "msg", "action", "spawn-sh", "--", cmd)

	return niriMsg.Run()
}

func NewRunner() Runner {
	return Runner{}
}
