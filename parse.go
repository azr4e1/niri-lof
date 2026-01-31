package nirilof

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	WindowIDConst          = "Window ID"
	IsFocusedConst         = "(focused)"
	TitleConst             = "Title"
	AppIDConst             = "App ID"
	IsFloatingConst        = "Is floating"
	PIDConst               = "PID"
	WorkspaceIDConst       = "Workspace ID"
	LayoutConst            = "Layout"
	TileSizeConst          = "Tile size"
	ScrollingPositionConst = "Scrolling position"
	WindowSizeConst        = "Window size"
	WindowOffsetConst      = "Window offset in tile"
	WorkSapceViewPosition  = "Workspace-view position"
	WindowBlockSeparator   = "\n\n"
)

func ParseNiriWindows(content string) ([]Window, error) {
	blocks := strings.Split(content, WindowBlockSeparator)
	cleanBlocks := []string{}
	for _, b := range blocks {
		block := strings.TrimSpace(b)
		if len(b) == 0 {
			continue
		}
		cleanBlocks = append(cleanBlocks, block)
	}

	windows := []Window{}
	for _, b := range cleanBlocks {
		window, err := ParseWindow(b)
		if err != nil {
			return nil, err
		}
		windows = append(windows, window)
	}

	return windows, nil
}

func ParseWindow(content string) (Window, error) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return Window{}, errors.New("empty container")
	}
	windowId, isFocused, err := parseWindowID(lines[0])
	if err != nil {
		return Window{}, err
	}

	// base indentation should be 2
	valueMap, _ := parseIndentation(lines[1:], 2)

	window := Window{
		ID:      windowId,
		Focused: isFocused,
	}
	for key, val := range valueMap {
		switch v := val.(type) {
		case string:
			switch key {
			case TitleConst:
				window.Title = strings.Trim(v, "\"")
			case AppIDConst:
				window.AppID = strings.Trim(v, "\"")
			case IsFloatingConst:
				var isFloating bool
				if val == "yes" {
					isFloating = true
				}
				window.IsFloating = isFloating
			case PIDConst:
				pID, err := strconv.Atoi(strings.TrimSpace(v))
				if err != nil {
					return Window{}, fmt.Errorf("couldn't parse PID: %w", err)
				}
				window.PID = pID
			case WorkspaceIDConst:
				wID, err := strconv.Atoi(strings.TrimSpace(v))
				if err != nil {
					return Window{}, fmt.Errorf("couldn't parse Workspace ID: %w", err)
				}
				window.WorkspaceID = wID
			default:
				continue
			}
		case map[string]any:
			layout, err := parseLayout(v)
			if err != nil {
				return Window{}, fmt.Errorf("couldn't parse Layout: %w", err)
			}
			window.Layout = layout
		default:
			continue
		}

	}
	return window, nil
}

func parseLayout(layoutMap map[string]any) (Layout, error) {
	layout := Layout{}
	for key, val := range layoutMap {
		switch v := val.(type) {
		case string:
			switch key {
			case TileSizeConst:
				size, err := parseSize(v)
				if err != nil {
					return Layout{}, err
				}
				layout.TileSize = size
			case WindowSizeConst:
				size, err := parseSize(v)
				if err != nil {
					return Layout{}, err
				}
				layout.WindowSize = size
			case WindowOffsetConst:
				size, err := parseSize(v)
				if err != nil {
					return Layout{}, err
				}
				layout.WindowOffsetTile = size
			case ScrollingPositionConst:
				position, err := parsePosition(v)
				if err != nil {
					return Layout{}, err
				}
				layout.ScrollingPos = position
			case WorkSapceViewPosition:
				position, err := parseFloatingPosition(v)
				if err != nil {
					return Layout{}, err
				}
				layout.WorkspaceViewPosition = position
			}
		default:
			continue
		}
	}

	return layout, nil
}

func parseSize(value string) (Size, error) {
	parts := strings.Split(value, "x")
	if len(parts) != 2 {
		return Size{}, fmt.Errorf("error parsing size: %s", value)
	}
	width, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return Size{}, fmt.Errorf("error parsing width size: %w", err)
	}
	height, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Size{}, fmt.Errorf("error parsing height size: %w", err)
	}
	return Size{width, height}, nil
}

func parsePosition(value string) (Position, error) {
	parts := strings.Split(value, ",")
	if len(parts) != 2 {
		return Position{}, fmt.Errorf("error parsing floating position: %s", value)
	}
	columnVal := strings.TrimPrefix(strings.TrimSpace(parts[0]), "column")
	column, err := strconv.Atoi(strings.TrimSpace(columnVal))
	if err != nil {
		return Position{}, fmt.Errorf("error parsing column position: %w", err)
	}
	tileVal := strings.TrimPrefix(strings.TrimSpace(parts[1]), "tile")
	tile, err := strconv.Atoi(strings.TrimSpace(tileVal))
	if err != nil {
		return Position{}, fmt.Errorf("error parsing tile position: %w", err)
	}
	return Position{column, tile}, nil
}

func parseFloatingPosition(value string) (FloatingPosition, error) {
	parts := strings.Split(value, ",")
	if len(parts) != 2 {
		return FloatingPosition{}, fmt.Errorf("error parsing floating position: %s", value)
	}
	width, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return FloatingPosition{}, fmt.Errorf("error parsing x position: %w", err)
	}
	height, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return FloatingPosition{}, fmt.Errorf("error parsing y position: %w", err)
	}
	return FloatingPosition{width, height}, nil
}

func getSpaceIndentation(line string) int {
	var level int
	for ; level < len(line) && line[level] == ' '; level++ {
	}

	return level
}

func parseIndentation(lines []string, currIndentationLevel int) (map[string]any, int) {
	if len(lines) == 0 {
		return nil, 0
	}
	parsedLines := map[string]any{}
	var lastVal string
	var i int
	for i = 0; i < len(lines); i++ {
		indentationLevel := getSpaceIndentation(lines[i])
		if indentationLevel > currIndentationLevel {
			newParsedLines, j := parseIndentation(lines[i:], indentationLevel)
			i += j
			parsedLines[lastVal] = newParsedLines
		} else if indentationLevel < currIndentationLevel {
			return parsedLines, i
		} else {
			vals := strings.SplitN(lines[i], ":", 2)
			if len(vals) != 2 {
				continue
			}
			trimmedName := strings.TrimSpace(vals[0])
			trimmedVal := strings.TrimSpace(vals[1])
			parsedLines[trimmedName] = trimmedVal
			lastVal = trimmedName
		}
	}
	return parsedLines, i
}

func parseWindowID(line string) (int, bool, error) {
	if !strings.HasPrefix(line, WindowIDConst) {
		return 0, false, errors.New("doesn't match window type")
	}
	var id int
	var isFocused bool
	parts := strings.Split(line, ":")
	if len(parts) <= 1 {
		return 0, false, errors.New("doesn't match window type")
	}
	if second := parts[1]; len(second) != 0 && strings.HasSuffix(second, IsFocusedConst) {
		isFocused = true
	}
	idString := strings.TrimSpace(strings.TrimPrefix(parts[0], WindowIDConst))
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, false, err
	}

	return id, isFocused, nil
}
