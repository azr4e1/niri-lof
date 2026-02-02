package cmdline

import (
	"fmt"
	"strconv"
	"strings"

	nirilof "github.com/azr4e1/niri-lof"
)

const FormatSep = ","

const (
	Focused     = "FOCUSED"
	ID          = "ID"
	Title       = "TITLE"
	AppID       = "APPID"
	IsFloating  = "ISFLOATING"
	PID         = "PID"
	WorkspaceID = "WORKSPACEID"
)

// Format tableOptions
var tableOptions = map[string]int{
	ID:          0,
	AppID:       1,
	Title:       2,
	Focused:     3,
	IsFloating:  4,
	PID:         5,
	WorkspaceID: 6,
}

func ValidateFormat(format string) ([]string, error) {
	if strings.TrimSpace(format) == "" {
		return []string{}, nil
	}
	formatOptions := strings.Split(strings.TrimSpace(format), FormatSep)
	tableFormat := []string{}
	for _, opt := range formatOptions {
		_, ok := tableOptions[opt]
		if !ok {
			return nil, fmt.Errorf("'%s' is not a format option", opt)
		}
		tableFormat = append(tableFormat, opt)
	}
	return tableFormat, nil
}

func FormatWindows(windows []nirilof.Window, options []string) string {
	if len(options) == 0 {
		options = make([]string, len(tableOptions))
		for k, i := range tableOptions {
			options[i] = k
		}
	}

	rows := [][]string{}

	// add the header
	rows = append(rows, options)

	// add the values
	for _, w := range windows {
		newRow := []string{}
		for _, c := range options {
			switch c {
			case ID:
				newRow = append(newRow, fmt.Sprintf("%d", w.ID))
			case AppID:
				newRow = append(newRow, w.AppID)
			case Title:
				newRow = append(newRow, w.Title)
			case Focused:
				newRow = append(newRow, fmt.Sprintf("%t", w.Focused))
			case IsFloating:
				newRow = append(newRow, fmt.Sprintf("%t", w.IsFloating))
			case PID:
				newRow = append(newRow, fmt.Sprintf("%d", w.PID))
			case WorkspaceID:
				newRow = append(newRow, fmt.Sprintf("%d", w.WorkspaceID))
			default:
				continue
			}
		}
		rows = append(rows, newRow)
	}

	// for each header, get the length of the longest value
	longest := map[string]int{}
	for _, row := range rows {
		for i, c := range row {
			header := options[i]
			val, ok := longest[header]
			if !ok {
				longest[header] = val
			}
			if lenC := len(c); lenC > val {
				longest[header] = lenC
			}
		}
	}

	// format rows
	rowStrings := []string{}
	for _, row := range rows {
		cols := []string{}
		for i, c := range row {
			header := options[i]
			length := longest[header]
			strLength := strconv.Itoa(length)
			valString := fmt.Sprintf("%-"+strLength+"s", c)
			cols = append(cols, valString)
		}
		newRow := strings.Join(cols, "  ")
		rowStrings = append(rowStrings, newRow)
	}
	table := strings.Join(rowStrings, "\n")

	return table
}
