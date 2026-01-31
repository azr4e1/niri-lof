package nirilof

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/google/shlex"
)

func GetWindows() ([]Window, error) {
	niriMsg := exec.Command("niri", "msg", "windows")
	data, err := niriMsg.CombinedOutput()
	if err != nil {
		return nil, err
	}

	windows, err := ParseNiriWindows(string(data))

	return windows, err
}

func FindWindowByAppID(appID string, windows []Window) []Window {
	appIDWindows := []Window{}

	for _, w := range windows {
		if w.AppID == appID {
			appIDWindows = append(appIDWindows, w)
		}
	}

	return appIDWindows
}

func FocusWindow(window Window, allWindows []Window) error {
	var exists bool
	for _, w := range allWindows {
		if w.ID == window.ID {
			exists = true
			break
		}
	}

	if !exists {
		return fmt.Errorf("windows with ID %d does not exist\n", window.ID)
	}
	niriMsg := exec.Command("niri", "msg", "action", "focus-window", "--id", fmt.Sprintf("%d", window.ID))
	err := niriMsg.Run()

	return err
}

func ParseCommand(cmd string) (*exec.Cmd, error) {
	shellSplit, err := shlex.Split(cmd)
	if err != nil {
		return nil, err
	}
	if len(shellSplit) == 0 {
		return nil, errors.New("empty command string")
	}
	name := shellSplit[0]
	args := shellSplit[1:]

	return exec.Command(name, args...), nil
}

func LaunchOrFocus(appID string, cmd string) error {
	command, err := ParseCommand(cmd)
	if err != nil {
		return err
	}
	allWindows, err := GetWindows()
	if err != nil {
		return err
	}

	appIDWindows := FindWindowByAppID(appID, allWindows)
	if len(appIDWindows) == 0 {
		err := command.Run()
		return err
	}

	window := appIDWindows[0]
	return FocusWindow(window, allWindows)
}
