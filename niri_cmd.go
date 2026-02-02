package nirilof

import (
	"fmt"
	"slices"
	"strings"
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
	if err != nil {
		return nil, err
	}

	// sort them by ID
	slices.SortFunc(windows, func(w1, w2 Window) int {
		switch {
		case w1.ID < w2.ID:
			return -1
		case w1.ID == w2.ID:
			return 0
		default:
			return 1
		}
	})

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

// Given a list of windows, find the index of the focused window.
// If none focused, returns -1
func FindFocusedWindow(allWindows []Window) int {
	index := -1
	for i, w := range allWindows {
		if w.Focused {
			index = i
			break
		}
	}

	return index
}

// Get next windows from index provided
func GetNextWindow(currentIndex int, allWindows []Window) Window {
	if len(allWindows) == 0 {
		return Window{}
	}

	if currentIndex < 0 || currentIndex >= len(allWindows)-1 {
		return allWindows[0]
	}

	return allWindows[currentIndex+1]
}

// Find a window by appID. If there are windows with that appID, focus the window
// with next ID of the currently focused, if the currently focused also has same appID.
// Otherwise, run the command cmd provided
func LaunchOrFocus(runner NiriRunner, appID string, cmd string) error {
	allWindows, err := GetWindows(runner)
	if err != nil {
		return err
	}

	appIDWindows := FindWindowByAppID(appID, allWindows)
	if len(appIDWindows) == 0 {
		// if cmd is null, ignore
		if len(strings.TrimSpace(cmd)) == 0 {
			return nil
		}
		err = runner.Spawn(cmd)
		return err
	}

	windowIndex := FindFocusedWindow(appIDWindows)
	window := GetNextWindow(windowIndex, appIDWindows)

	return FocusWindow(runner, window, allWindows)
}
