package nirilof

import (
	"reflect"
	"testing"
)

func TestParseWindowID(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantID    int
		wantFocus bool
		wantErr   bool
	}{
		{
			name:      "unfocused window",
			line:      "Window ID 7:",
			wantID:    7,
			wantFocus: false,
		},
		{
			name:      "focused window",
			line:      "Window ID 15: (focused)",
			wantID:    15,
			wantFocus: true,
		},
		{
			name:      "large ID",
			line:      "Window ID 99999:",
			wantID:    99999,
			wantFocus: false,
		},
		{
			name:    "missing prefix",
			line:    "Something else 7:",
			wantErr: true,
		},
		{
			name:    "no colon separator",
			line:    "Window ID 7",
			wantErr: true,
		},
		{
			name:    "non-numeric ID",
			line:    "Window ID abc:",
			wantErr: true,
		},
		{
			name:    "empty string",
			line:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, focused, err := parseWindowID(tt.line)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if id != tt.wantID {
				t.Errorf("ID = %d, want %d", id, tt.wantID)
			}
			if focused != tt.wantFocus {
				t.Errorf("focused = %v, want %v", focused, tt.wantFocus)
			}
		})
	}
}

func TestGetSpaceIndentation(t *testing.T) {
	tests := []struct {
		name string
		line string
		want int
	}{
		{"no indentation", "hello", 0},
		{"two spaces", "  hello", 2},
		{"four spaces", "    hello", 4},
		{"only spaces", "    ", 4},
		{"empty string", "", 0},
		{"tab is not a space", "\thello", 0},
		{"mixed tab after spaces", "  \thello", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSpaceIndentation(tt.line)
			if got != tt.want {
				t.Errorf("getSpaceIndentation(%q) = %d, want %d", tt.line, got, tt.want)
			}
		})
	}
}

func TestParseSize(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    NumericalPair
		wantErr bool
	}{
		{
			name:  "normal size",
			value: "1491 x 1060",
			want:  NumericalPair{1491, 1060},
		},
		{
			name:  "zero size",
			value: "0 x 0",
			want:  NumericalPair{0, 0},
		},
		{
			name:  "no spaces around x",
			value: "100x200",
			want:  NumericalPair{100, 200},
		},
		{
			name:    "missing separator",
			value:   "1491 1060",
			wantErr: true,
		},
		{
			name:    "too many parts",
			value:   "1 x 2 x 3",
			wantErr: true,
		},
		{
			name:    "non-numeric width",
			value:   "abc x 100",
			wantErr: true,
		},
		{
			name:    "non-numeric height",
			value:   "100 x abc",
			wantErr: true,
		},
		{
			name:    "empty string",
			value:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSize(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSize(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestParsePosition(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    NumericalPair
		wantErr bool
	}{
		{
			name:  "normal position",
			value: "column 1, tile 1",
			want:  NumericalPair{1, 1},
		},
		{
			name:  "larger values",
			value: "column 5, tile 3",
			want:  NumericalPair{5, 3},
		},
		{
			name:    "missing comma",
			value:   "column 1 tile 1",
			wantErr: true,
		},
		{
			name:    "non-numeric column",
			value:   "column abc, tile 1",
			wantErr: true,
		},
		{
			name:    "non-numeric tile",
			value:   "column 1, tile abc",
			wantErr: true,
		},
		{
			name:    "empty string",
			value:   "",
			wantErr: true,
		},
		{
			name:    "too many commas",
			value:   "column 1, tile 1, extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePosition(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePosition(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestParseFloatingPosition(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    NumericalPair
		wantErr bool
	}{
		{
			name:  "normal position",
			value: "53, 20",
			want:  NumericalPair{53, 20},
		},
		{
			name:  "zero position",
			value: "0, 0",
			want:  NumericalPair{0, 0},
		},
		{
			name:  "no spaces",
			value: "10,20",
			want:  NumericalPair{10, 20},
		},
		{
			name:    "missing comma",
			value:   "53 20",
			wantErr: true,
		},
		{
			name:    "non-numeric x",
			value:   "abc, 20",
			wantErr: true,
		},
		{
			name:    "non-numeric y",
			value:   "53, abc",
			wantErr: true,
		},
		{
			name:    "empty string",
			value:   "",
			wantErr: true,
		},
		{
			name:    "too many commas",
			value:   "1, 2, 3",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFloatingPosition(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFloatingPosition(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestParseIndentation(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		result, consumed := parseIndentation([]string{}, 2)
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
		if consumed != 0 {
			t.Errorf("consumed = %d, want 0", consumed)
		}
	})

	t.Run("flat key-value pairs", func(t *testing.T) {
		lines := []string{
			`  Title: "Firefox"`,
			`  App ID: "firefox"`,
			`  PID: 1234`,
		}
		result, consumed := parseIndentation(lines, 2)
		if consumed != 3 {
			t.Errorf("consumed = %d, want 3", consumed)
		}
		if result["Title"] != `"Firefox"` {
			t.Errorf("Title = %v, want %q", result["Title"], `"Firefox"`)
		}
		if result["App ID"] != `"firefox"` {
			t.Errorf("App ID = %v, want %q", result["App ID"], `"firefox"`)
		}
		if result["PID"] != "1234" {
			t.Errorf("PID = %v, want %q", result["PID"], "1234")
		}
	})

	t.Run("nested structure", func(t *testing.T) {
		lines := []string{
			"  Layout:",
			"    Tile size: 1491 x 1060",
			"    Window size: 1491 x 1060",
		}
		result, _ := parseIndentation(lines, 2)
		nested, ok := result["Layout"].(map[string]any)
		if !ok {
			t.Fatalf("Layout should be map[string]any, got %T", result["Layout"])
		}
		if nested["Tile size"] != "1491 x 1060" {
			t.Errorf("Tile size = %v, want %q", nested["Tile size"], "1491 x 1060")
		}
		if nested["Window size"] != "1491 x 1060" {
			t.Errorf("Window size = %v, want %q", nested["Window size"], "1491 x 1060")
		}
	})

	t.Run("returns on decreased indentation", func(t *testing.T) {
		lines := []string{
			"    Tile size: 100 x 200",
			"  Next field: value",
		}
		result, consumed := parseIndentation(lines, 4)
		if consumed != 1 {
			t.Errorf("consumed = %d, want 1", consumed)
		}
		if result["Tile size"] != "100 x 200" {
			t.Errorf("Tile size = %v, want %q", result["Tile size"], "100 x 200")
		}
	})

	t.Run("multi-level nesting with returns to each level", func(t *testing.T) {
		// Structure:
		//   A: val_a                  (baseline level 2)
		//   B:                        (baseline level 2, triggers nesting)
		//     C: val_c               (level 4)
		//     D:                     (level 4, triggers deeper nesting)
		//       E: val_e             (level 6)
		//     F: val_f               (back to level 4)
		//   G: val_g                  (back to baseline level 2)
		lines := []string{
			"  A: val_a",
			"  B:",
			"    C: val_c",
			"    D:",
			"      E: val_e",
			"    F: val_f",
			"  G: val_g",
		}
		result, consumed := parseIndentation(lines, 2)
		if consumed != 7 {
			t.Errorf("consumed = %d, want 7", consumed)
		}

		// Baseline keys
		if result["A"] != "val_a" {
			t.Errorf("A = %v, want %q", result["A"], "val_a")
		}
		if result["G"] != "val_g" {
			t.Errorf("G = %v, want %q", result["G"], "val_g")
		}

		// First nesting: B -> {C, D, F}
		bMap, ok := result["B"].(map[string]any)
		if !ok {
			t.Fatalf("B should be map[string]any, got %T", result["B"])
		}
		if bMap["C"] != "val_c" {
			t.Errorf("B.C = %v, want %q", bMap["C"], "val_c")
		}
		if bMap["F"] != "val_f" {
			t.Errorf("B.F = %v, want %q", bMap["F"], "val_f")
		}

		// Second nesting: D -> {E}
		dMap, ok := bMap["D"].(map[string]any)
		if !ok {
			t.Fatalf("B.D should be map[string]any, got %T", bMap["D"])
		}
		if dMap["E"] != "val_e" {
			t.Errorf("B.D.E = %v, want %q", dMap["E"], "val_e")
		}
	})

	t.Run("line without colon is skipped", func(t *testing.T) {
		lines := []string{
			"  valid: yes",
			"  no-colon-here",
			"  also valid: no",
		}
		result, consumed := parseIndentation(lines, 2)
		if consumed != 3 {
			t.Errorf("consumed = %d, want 3", consumed)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 entries, got %d", len(result))
		}
	})
}

func TestParseLayout(t *testing.T) {
	t.Run("tiled window layout", func(t *testing.T) {
		layoutMap := map[string]any{
			"Tile size":             "1491 x 1060",
			"Scrolling position":    "column 1, tile 1",
			"Window size":           "1491 x 1060",
			"Window offset in tile": "0 x 0",
		}
		got, err := parseLayout(layoutMap)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := Layout{
			TileSize:         NumericalPair{1491, 1060},
			ScrollingPos:     NumericalPair{1, 1},
			WindowSize:       NumericalPair{1491, 1060},
			WindowOffsetTile: NumericalPair{0, 0},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("floating window layout", func(t *testing.T) {
		layoutMap := map[string]any{
			"Tile size":               "1867 x 1060",
			"Workspace-view position": "53, 20",
			"Window size":             "1867 x 1060",
			"Window offset in tile":   "0 x 0",
		}
		got, err := parseLayout(layoutMap)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := Layout{
			TileSize:              NumericalPair{1867, 1060},
			WorkspaceViewPosition: NumericalPair{53, 20},
			WindowSize:            NumericalPair{1867, 1060},
			WindowOffsetTile:      NumericalPair{0, 0},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("bad tile size", func(t *testing.T) {
		layoutMap := map[string]any{
			"Tile size": "bad",
		}
		_, err := parseLayout(layoutMap)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("bad window size", func(t *testing.T) {
		layoutMap := map[string]any{
			"Window size": "bad",
		}
		_, err := parseLayout(layoutMap)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("bad window offset", func(t *testing.T) {
		layoutMap := map[string]any{
			"Window offset in tile": "bad",
		}
		_, err := parseLayout(layoutMap)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("bad scrolling position", func(t *testing.T) {
		layoutMap := map[string]any{
			"Scrolling position": "bad",
		}
		_, err := parseLayout(layoutMap)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("bad workspace-view position", func(t *testing.T) {
		layoutMap := map[string]any{
			"Workspace-view position": "bad",
		}
		_, err := parseLayout(layoutMap)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("non-string values are skipped", func(t *testing.T) {
		layoutMap := map[string]any{
			"Tile size":    "100 x 200",
			"unknown-nest": map[string]any{"a": "b"},
		}
		got, err := parseLayout(layoutMap)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got.TileSize, NumericalPair{100, 200}) {
			t.Errorf("TileSize = %v, want {100 200}", got.TileSize)
		}
	})

	t.Run("empty map", func(t *testing.T) {
		got, err := parseLayout(map[string]any{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, Layout{}) {
			t.Errorf("expected zero Layout, got %+v", got)
		}
	})
}

func TestParseWindow(t *testing.T) {
	t.Run("tiled unfocused window", func(t *testing.T) {
		content := `Window ID 7:
  Title: "A Tour of Go — Mozilla Firefox"
  App ID: "firefox"
  Is floating: no
  PID: 3798
  Workspace ID: 4
  Layout:
    Tile size: 1491 x 1060
    Scrolling position: column 1, tile 1
    Window size: 1491 x 1060
    Window offset in tile: 0 x 0`

		got, err := ParseWindow(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != 7 {
			t.Errorf("ID = %d, want 7", got.ID)
		}
		if got.Focused {
			t.Error("expected Focused = false")
		}
		if got.Title != "A Tour of Go — Mozilla Firefox" {
			t.Errorf("Title = %q, want %q", got.Title, "A Tour of Go — Mozilla Firefox")
		}
		if got.AppID != "firefox" {
			t.Errorf("AppID = %q, want %q", got.AppID, "firefox")
		}
		if got.IsFloating {
			t.Error("expected IsFloating = false")
		}
		if got.PID != 3798 {
			t.Errorf("PID = %d, want 3798", got.PID)
		}
		if got.WorkspaceID != 4 {
			t.Errorf("WorkspaceID = %d, want 4", got.WorkspaceID)
		}
		if !reflect.DeepEqual(got.Layout.TileSize, NumericalPair{1491, 1060}) {
			t.Errorf("TileSize = %v, want {1491 1060}", got.Layout.TileSize)
		}
		if !reflect.DeepEqual(got.Layout.ScrollingPos, NumericalPair{1, 1}) {
			t.Errorf("ScrollingPos = %v, want {1 1}", got.Layout.ScrollingPos)
		}
		if !reflect.DeepEqual(got.Layout.WindowSize, NumericalPair{1491, 1060}) {
			t.Errorf("WindowSize = %v, want {1491 1060}", got.Layout.WindowSize)
		}
		if !reflect.DeepEqual(got.Layout.WindowOffsetTile, NumericalPair{0, 0}) {
			t.Errorf("WindowOffsetTile = %v, want {0 0}", got.Layout.WindowOffsetTile)
		}
	})

	t.Run("floating focused window", func(t *testing.T) {
		content := `Window ID 15: (focused)
  Title: "ld@archbox:~/Desktop/Projects/niri-lof/cmd/main"
  App ID: "Alacritty"
  Is floating: yes
  PID: 9626
  Workspace ID: 5
  Layout:
    Tile size: 1867 x 1060
    Workspace-view position: 53, 20
    Window size: 1867 x 1060
    Window offset in tile: 0 x 0`

		got, err := ParseWindow(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != 15 {
			t.Errorf("ID = %d, want 15", got.ID)
		}
		if !got.Focused {
			t.Error("expected Focused = true")
		}
		if got.AppID != "Alacritty" {
			t.Errorf("AppID = %q, want %q", got.AppID, "Alacritty")
		}
		if !got.IsFloating {
			t.Error("expected IsFloating = true")
		}
		if got.PID != 9626 {
			t.Errorf("PID = %d, want 9626", got.PID)
		}
		if !reflect.DeepEqual(got.Layout.WorkspaceViewPosition, NumericalPair{53, 20}) {
			t.Errorf("WorkspaceViewPosition = %v, want {53 20}", got.Layout.WorkspaceViewPosition)
		}
	})

	t.Run("empty content", func(t *testing.T) {
		_, err := ParseWindow("")
		if err == nil {
			t.Fatal("expected error for empty content")
		}
	})

	t.Run("invalid first line", func(t *testing.T) {
		_, err := ParseWindow("Not a window line")
		if err == nil {
			t.Fatal("expected error for invalid first line")
		}
	})

	t.Run("invalid PID", func(t *testing.T) {
		content := `Window ID 1:
  Title: "test"
  App ID: "test"
  PID: notanumber
  Workspace ID: 1`

		_, err := ParseWindow(content)
		if err == nil {
			t.Fatal("expected error for non-numeric PID")
		}
	})

	t.Run("invalid workspace ID", func(t *testing.T) {
		content := `Window ID 1:
  Title: "test"
  App ID: "test"
  PID: 100
  Workspace ID: notanumber`

		_, err := ParseWindow(content)
		if err == nil {
			t.Fatal("expected error for non-numeric workspace ID")
		}
	})

	t.Run("invalid layout value", func(t *testing.T) {
		content := `Window ID 1:
  Title: "test"
  App ID: "test"
  PID: 100
  Workspace ID: 1
  Layout:
    Tile size: bad`

		_, err := ParseWindow(content)
		if err == nil {
			t.Fatal("expected error for bad layout")
		}
	})
}

func TestParseNiriWindows(t *testing.T) {
	t.Run("full test data", func(t *testing.T) {
		content := `Window ID 7:
  Title: "A Tour of Go — Original profile — Mozilla Firefox"
  App ID: "firefox"
  Is floating: no
  PID: 3798
  Workspace ID: 4
  Layout:
    Tile size: 1491 x 1060
    Scrolling position: column 1, tile 1
    Window size: 1491 x 1060
    Window offset in tile: 0 x 0

Window ID 8:
  Title: "Fwd: Action Required — Mozilla Thunderbird"
  App ID: "org.mozilla.Thunderbird"
  Is floating: no
  PID: 4441
  Workspace ID: 1
  Layout:
    Tile size: 772 x 701
    Scrolling position: column 1, tile 1
    Window size: 772 x 701
    Window offset in tile: 0 x 0

Window ID 15: (focused)
  Title: "ld@archbox:~/Desktop/Projects/niri-lof/cmd/main"
  App ID: "Alacritty"
  Is floating: yes
  PID: 9626
  Workspace ID: 5
  Layout:
    Tile size: 1867 x 1060
    Workspace-view position: 53, 20
    Window size: 1867 x 1060
    Window offset in tile: 0 x 0`

		windows, err := ParseNiriWindows(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(windows) != 3 {
			t.Fatalf("got %d windows, want 3", len(windows))
		}

		// First window
		if windows[0].ID != 7 {
			t.Errorf("windows[0].ID = %d, want 7", windows[0].ID)
		}
		if windows[0].AppID != "firefox" {
			t.Errorf("windows[0].AppID = %q, want %q", windows[0].AppID, "firefox")
		}
		if windows[0].Focused {
			t.Error("windows[0] should not be focused")
		}

		// Second window
		if windows[1].ID != 8 {
			t.Errorf("windows[1].ID = %d, want 8", windows[1].ID)
		}
		if windows[1].AppID != "org.mozilla.Thunderbird" {
			t.Errorf("windows[1].AppID = %q, want %q", windows[1].AppID, "org.mozilla.Thunderbird")
		}

		// Third window (focused, floating)
		if windows[2].ID != 15 {
			t.Errorf("windows[2].ID = %d, want 15", windows[2].ID)
		}
		if !windows[2].Focused {
			t.Error("windows[2] should be focused")
		}
		if !windows[2].IsFloating {
			t.Error("windows[2] should be floating")
		}
		if !reflect.DeepEqual(windows[2].Layout.WorkspaceViewPosition, NumericalPair{53, 20}) {
			t.Errorf("windows[2] workspace-view position = %v, want {53 20}",
				windows[2].Layout.WorkspaceViewPosition)
		}
	})

	t.Run("single window", func(t *testing.T) {
		content := `Window ID 1:
  Title: "test"
  App ID: "myapp"
  Is floating: no
  PID: 100
  Workspace ID: 1
  Layout:
    Tile size: 800 x 600
    Scrolling position: column 1, tile 1
    Window size: 800 x 600
    Window offset in tile: 0 x 0`

		windows, err := ParseNiriWindows(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(windows) != 1 {
			t.Fatalf("got %d windows, want 1", len(windows))
		}
		if windows[0].AppID != "myapp" {
			t.Errorf("AppID = %q, want %q", windows[0].AppID, "myapp")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		windows, err := ParseNiriWindows("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(windows) != 0 {
			t.Errorf("got %d windows, want 0", len(windows))
		}
	})

	t.Run("whitespace only", func(t *testing.T) {
		windows, err := ParseNiriWindows("   \n\n   \n  ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(windows) != 0 {
			t.Errorf("got %d windows, want 0", len(windows))
		}
	})

	t.Run("trailing newlines", func(t *testing.T) {
		content := `Window ID 1:
  Title: "test"
  App ID: "myapp"
  Is floating: no
  PID: 100
  Workspace ID: 1
  Layout:
    Tile size: 800 x 600
    Scrolling position: column 1, tile 1
    Window size: 800 x 600
    Window offset in tile: 0 x 0

`

		windows, err := ParseNiriWindows(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(windows) != 1 {
			t.Fatalf("got %d windows, want 1", len(windows))
		}
	})

	t.Run("invalid window block causes error", func(t *testing.T) {
		content := `Window ID 1:
  Title: "test"
  App ID: "myapp"
  PID: 100
  Workspace ID: 1

Not a valid window block`

		_, err := ParseNiriWindows(content)
		if err == nil {
			t.Fatal("expected error for invalid block")
		}
	})

	t.Run("all fields parsed from real data", func(t *testing.T) {
		content := `Window ID 7:
  Title: "A Tour of Go — Original profile — Mozilla Firefox"
  App ID: "firefox"
  Is floating: no
  PID: 3798
  Workspace ID: 4
  Layout:
    Tile size: 1491 x 1060
    Scrolling position: column 1, tile 1
    Window size: 1491 x 1060
    Window offset in tile: 0 x 0`

		windows, err := ParseNiriWindows(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := Window{
			ID:          7,
			Focused:     false,
			Title:       "A Tour of Go — Original profile — Mozilla Firefox",
			AppID:       "firefox",
			IsFloating:  false,
			PID:         3798,
			WorkspaceID: 4,
			Layout: Layout{
				TileSize:         NumericalPair{1491, 1060},
				ScrollingPos:     NumericalPair{1, 1},
				WindowSize:       NumericalPair{1491, 1060},
				WindowOffsetTile: NumericalPair{0, 0},
			},
		}

		if !reflect.DeepEqual(windows[0], want) {
			t.Errorf("got:\n%+v\nwant:\n%+v", windows[0], want)
		}
	})
}
