package nirilof

import (
	"testing"
)

// Sample JSON taken from real `niri msg -j windows` output.
var sampleJSON = []byte(`[
	{
		"id": 3,
		"title": "Inbox - Unified Folders - Mozilla Thunderbird",
		"app_id": "org.mozilla.Thunderbird",
		"pid": 95229,
		"workspace_id": 1,
		"is_focused": false,
		"is_floating": false,
		"is_urgent": false,
		"layout": {
			"pos_in_scrolling_layout": [1, 1],
			"tile_size": [772.0, 701.0],
			"window_size": [772, 701],
			"tile_pos_in_workspace_view": null,
			"window_offset_in_tile": [0.0, 0.0]
		},
		"focus_timestamp": {"secs": 31617, "nanos": 601451712}
	},
	{
		"id": 29,
		"title": "Terminal",
		"app_id": "Alacritty",
		"pid": 180761,
		"workspace_id": 3,
		"is_focused": true,
		"is_floating": false,
		"is_urgent": false,
		"layout": {
			"pos_in_scrolling_layout": [1, 1],
			"tile_size": [772.0, 701.0],
			"window_size": [772, 701],
			"tile_pos_in_workspace_view": null,
			"window_offset_in_tile": [0.0, 0.0]
		},
		"focus_timestamp": {"secs": 33518, "nanos": 498286808}
	},
	{
		"id": 36,
		"title": "Floating Terminal",
		"app_id": "Alacritty",
		"pid": 230702,
		"workspace_id": 4,
		"is_focused": false,
		"is_floating": true,
		"is_urgent": false,
		"layout": {
			"pos_in_scrolling_layout": null,
			"tile_size": [928.0, 1060.0],
			"window_size": [928, 1060],
			"tile_pos_in_workspace_view": [992.0, 20.0],
			"window_offset_in_tile": [0.0, 0.0]
		},
		"focus_timestamp": {"secs": 33511, "nanos": 985311335}
	}
]`)

func TestParseNiriWindowsJSON(t *testing.T) {
	t.Run("real niri output", func(t *testing.T) {
		windows, err := ParseNiriWindowsJSON(sampleJSON)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(windows) != 3 {
			t.Fatalf("expected 3 windows, got %d", len(windows))
		}

		// First window: Thunderbird
		w := windows[0]
		if w.ID != 3 {
			t.Errorf("window[0].ID = %d, want 3", w.ID)
		}
		if w.Title != "Inbox - Unified Folders - Mozilla Thunderbird" {
			t.Errorf("window[0].Title = %q, want %q", w.Title, "Inbox - Unified Folders - Mozilla Thunderbird")
		}
		if w.AppID != "org.mozilla.Thunderbird" {
			t.Errorf("window[0].AppID = %q, want %q", w.AppID, "org.mozilla.Thunderbird")
		}
		if w.PID != 95229 {
			t.Errorf("window[0].PID = %d, want 95229", w.PID)
		}
		if w.WorkspaceID != 1 {
			t.Errorf("window[0].WorkspaceID = %d, want 1", w.WorkspaceID)
		}
		if w.Focused {
			t.Error("window[0].Focused = true, want false")
		}
		if w.IsFloating {
			t.Error("window[0].IsFloating = true, want false")
		}

		// Second window: focused terminal
		w = windows[1]
		if w.ID != 29 {
			t.Errorf("window[1].ID = %d, want 29", w.ID)
		}
		if w.AppID != "Alacritty" {
			t.Errorf("window[1].AppID = %q, want %q", w.AppID, "Alacritty")
		}
		if !w.Focused {
			t.Error("window[1].Focused = false, want true")
		}

		// Third window: floating
		w = windows[2]
		if w.ID != 36 {
			t.Errorf("window[2].ID = %d, want 36", w.ID)
		}
		if !w.IsFloating {
			t.Error("window[2].IsFloating = false, want true")
		}
	})

	t.Run("layout fields parsed correctly", func(t *testing.T) {
		windows, err := ParseNiriWindowsJSON(sampleJSON)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Tiled window: has pos_in_scrolling_layout, null tile_pos_in_workspace_view
		layout := windows[0].Layout
		if len(layout.ScrollingPos) != 2 || layout.ScrollingPos[0] != 1 || layout.ScrollingPos[1] != 1 {
			t.Errorf("ScrollingPos = %v, want [1 1]", layout.ScrollingPos)
		}
		if len(layout.TileSize) != 2 || layout.TileSize[0] != 772 || layout.TileSize[1] != 701 {
			t.Errorf("TileSize = %v, want [772 701]", layout.TileSize)
		}
		if len(layout.WindowSize) != 2 || layout.WindowSize[0] != 772 || layout.WindowSize[1] != 701 {
			t.Errorf("WindowSize = %v, want [772 701]", layout.WindowSize)
		}
		if layout.WorkspaceViewPosition != nil {
			t.Errorf("WorkspaceViewPosition = %v, want nil (null in JSON)", layout.WorkspaceViewPosition)
		}
		if len(layout.WindowOffsetTile) != 2 || layout.WindowOffsetTile[0] != 0 || layout.WindowOffsetTile[1] != 0 {
			t.Errorf("WindowOffsetTile = %v, want [0 0]", layout.WindowOffsetTile)
		}

		// Floating window: null pos_in_scrolling_layout, has tile_pos_in_workspace_view
		layout = windows[2].Layout
		if layout.ScrollingPos != nil {
			t.Errorf("floating ScrollingPos = %v, want nil (null in JSON)", layout.ScrollingPos)
		}
		if len(layout.WorkspaceViewPosition) != 2 || layout.WorkspaceViewPosition[0] != 992 || layout.WorkspaceViewPosition[1] != 20 {
			t.Errorf("WorkspaceViewPosition = %v, want [992 20]", layout.WorkspaceViewPosition)
		}
	})

	t.Run("empty array", func(t *testing.T) {
		windows, err := ParseNiriWindowsJSON([]byte(`[]`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(windows) != 0 {
			t.Errorf("expected 0 windows, got %d", len(windows))
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := ParseNiriWindowsJSON([]byte(`not json`))
		if err == nil {
			t.Fatal("expected error for invalid JSON, got nil")
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		_, err := ParseNiriWindowsJSON([]byte(`[{"id": 1, "title":}]`))
		if err == nil {
			t.Fatal("expected error for malformed JSON, got nil")
		}
	})

	t.Run("empty input", func(t *testing.T) {
		_, err := ParseNiriWindowsJSON([]byte(``))
		if err == nil {
			t.Fatal("expected error for empty input, got nil")
		}
	})

	t.Run("extra fields are ignored", func(t *testing.T) {
		input := []byte(`[{"id": 5, "app_id": "test", "is_urgent": false, "focus_timestamp": {"secs": 100, "nanos": 200}, "some_future_field": true}]`)
		windows, err := ParseNiriWindowsJSON(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(windows) != 1 {
			t.Fatalf("expected 1 window, got %d", len(windows))
		}
		if windows[0].ID != 5 {
			t.Errorf("ID = %d, want 5", windows[0].ID)
		}
	})

	t.Run("missing fields default to zero values", func(t *testing.T) {
		input := []byte(`[{"id": 10, "app_id": "minimal"}]`)
		windows, err := ParseNiriWindowsJSON(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(windows) != 1 {
			t.Fatalf("expected 1 window, got %d", len(windows))
		}
		w := windows[0]
		if w.Title != "" {
			t.Errorf("Title = %q, want empty string", w.Title)
		}
		if w.Focused {
			t.Error("Focused = true, want false")
		}
		if w.PID != 0 {
			t.Errorf("PID = %d, want 0", w.PID)
		}
	})

	t.Run("object instead of array", func(t *testing.T) {
		_, err := ParseNiriWindowsJSON([]byte(`{"id": 1}`))
		if err == nil {
			t.Fatal("expected error for JSON object instead of array, got nil")
		}
	})

	t.Run("null JSON", func(t *testing.T) {
		windows, err := ParseNiriWindowsJSON([]byte(`null`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(windows) != 0 {
			t.Errorf("expected 0 windows, got %d", len(windows))
		}
	})
}
