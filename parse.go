package nirilof

import (
	"encoding/json"
)

func ParseNiriWindowsJSON(jsonContent []byte) ([]Window, error) {
	var windows []Window
	err := json.Unmarshal(jsonContent, &windows)
	if err != nil {
		return windows, err
	}

	return windows, nil
}
