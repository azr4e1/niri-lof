package nirilof

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/google/shlex"
)

// Get all currently open windows in niri
func GetWindows() ([]Window, error) {
	niriMsg := exec.Command("niri", "msg", "-j", "windows")
	data, err := niriMsg.Output()
	if err != nil {
		return nil, err
	}

	windows, err := ParseNiriWindowsJSON(data)

	return windows, err
}

// Find a window by App ID among all the open windows in niri
func FindWindowByAppID(appID string, windows []Window) []Window {
	appIDWindows := []Window{}

	for _, w := range windows {
		if w.AppID == appID {
			appIDWindows = append(appIDWindows, w)
		}
	}

	return appIDWindows
}

// Focus a window in niri
func FocusWindow(window Window, allWindows []Window) error {
	var exists bool
	for _, w := range allWindows {
		if w.ID == window.ID {
			exists = true
			break
		}
	}

	if !exists {
		return fmt.Errorf("window with ID %d does not exist", window.ID)
	}
	niriMsg := exec.Command("niri", "msg", "action", "focus-window", "--id", fmt.Sprintf("%d", window.ID))
	err := niriMsg.Run()

	return err
}

// use shlex to split a string according to shell
// rules and create a command
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

// Find a window by appID. If there is, focus the first result.
// Otherwise, run the command cmd provided
func LaunchOrFocus(appID string, cmd string) error {
	allWindows, err := GetWindows()
	if err != nil {
		return err
	}

	appIDWindows := FindWindowByAppID(appID, allWindows)
	if len(appIDWindows) == 0 {
		command, err := ParseCommand(cmd)
		if err != nil {
			return err
		}

		err = command.Run()
		return err
	}

	window := appIDWindows[0]
	return FocusWindow(window, allWindows)
}
