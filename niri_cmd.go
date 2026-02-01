package nirilof

import (
	"fmt"
)

type NiriRunner interface {
	GetJSON() ([]byte, error)
	Focus(winID int) error
	Spawn(cmd string) error
}

// Get all currently open windows in niri
func GetWindows(runner NiriRunner) ([]Window, error) {
	data, err := runner.GetJSON()
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
func FocusWindow(runner NiriRunner, window Window, allWindows []Window) error {
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
	err := runner.Focus(window.ID)

	return err
}

// Find a window by appID. If there is, focus the first result.
// Otherwise, run the command cmd provided
func LaunchOrFocus(runner NiriRunner, appID string, cmd string) error {
	allWindows, err := GetWindows(runner)
	if err != nil {
		return err
	}

	appIDWindows := FindWindowByAppID(appID, allWindows)
	if len(appIDWindows) == 0 {
		err = runner.Spawn(cmd)
		return err
	}

	window := appIDWindows[0]
	return FocusWindow(runner, window, allWindows)
}
