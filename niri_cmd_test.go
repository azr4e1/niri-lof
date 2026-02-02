package nirilof

import (
	"errors"
	"testing"
)

// mockRunner implements NiriRunner for testing.
type mockRunner struct {
	jsonData  []byte
	jsonErr   error
	focusedID int
	focusErr  error
	spawnedCmd string
	spawnErr  error
}

func (m *mockRunner) GetJSON() ([]byte, error) {
	return m.jsonData, m.jsonErr
}

func (m *mockRunner) Focus(windowID int) error {
	m.focusedID = windowID
	return m.focusErr
}

func (m *mockRunner) Spawn(cmd string) error {
	m.spawnedCmd = cmd
	return m.spawnErr
}

var testWindows = []Window{
	{ID: 1, Title: "Firefox", AppID: "firefox", PID: 100, WorkspaceID: 1, Focused: true},
	{ID: 2, Title: "Terminal", AppID: "Alacritty", PID: 200, WorkspaceID: 2},
	{ID: 3, Title: "Another Terminal", AppID: "Alacritty", PID: 300, WorkspaceID: 3},
}

var testWindowsJSON = []byte(`[
	{"id":1,"title":"Firefox","app_id":"firefox","pid":100,"workspace_id":1,"is_focused":true,"is_floating":false},
	{"id":2,"title":"Terminal","app_id":"Alacritty","pid":200,"workspace_id":2,"is_focused":false,"is_floating":false},
	{"id":3,"title":"Another Terminal","app_id":"Alacritty","pid":300,"workspace_id":3,"is_focused":false,"is_floating":false}
]`)

func TestGetWindows(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		runner := &mockRunner{jsonData: testWindowsJSON}
		windows, err := GetWindows(runner)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(windows) != 3 {
			t.Fatalf("expected 3 windows, got %d", len(windows))
		}
		if windows[0].AppID != "firefox" {
			t.Errorf("windows[0].AppID = %q, want %q", windows[0].AppID, "firefox")
		}
	})

	t.Run("runner error", func(t *testing.T) {
		runner := &mockRunner{jsonErr: errors.New("niri not running")}
		_, err := GetWindows(runner)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		runner := &mockRunner{jsonData: []byte(`not json`)}
		_, err := GetWindows(runner)
		if err == nil {
			t.Fatal("expected error for invalid JSON, got nil")
		}
	})
}

func TestFindWindowByAppID(t *testing.T) {
	t.Run("single match", func(t *testing.T) {
		result := FindWindowByAppID("firefox", testWindows)
		if len(result) != 1 {
			t.Fatalf("expected 1 match, got %d", len(result))
		}
		if result[0].ID != 1 {
			t.Errorf("ID = %d, want 1", result[0].ID)
		}
	})

	t.Run("multiple matches", func(t *testing.T) {
		result := FindWindowByAppID("Alacritty", testWindows)
		if len(result) != 2 {
			t.Fatalf("expected 2 matches, got %d", len(result))
		}
		if result[0].ID != 2 || result[1].ID != 3 {
			t.Errorf("got IDs %d and %d, want 2 and 3", result[0].ID, result[1].ID)
		}
	})

	t.Run("no match", func(t *testing.T) {
		result := FindWindowByAppID("nonexistent", testWindows)
		if len(result) != 0 {
			t.Errorf("expected 0 matches, got %d", len(result))
		}
	})

	t.Run("empty window list", func(t *testing.T) {
		result := FindWindowByAppID("firefox", []Window{})
		if len(result) != 0 {
			t.Errorf("expected 0 matches, got %d", len(result))
		}
	})

	t.Run("case sensitive", func(t *testing.T) {
		result := FindWindowByAppID("Firefox", testWindows)
		if len(result) != 0 {
			t.Errorf("expected 0 matches (case mismatch), got %d", len(result))
		}
	})
}

func TestFocusWindow(t *testing.T) {
	t.Run("focuses existing window", func(t *testing.T) {
		runner := &mockRunner{}
		err := FocusWindow(runner, testWindows[1], testWindows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if runner.focusedID != 2 {
			t.Errorf("focused ID = %d, want 2", runner.focusedID)
		}
	})

	t.Run("window not in list", func(t *testing.T) {
		runner := &mockRunner{}
		missing := Window{ID: 999, AppID: "ghost"}
		err := FocusWindow(runner, missing, testWindows)
		if err == nil {
			t.Fatal("expected error for missing window, got nil")
		}
		if runner.focusedID != 0 {
			t.Errorf("Focus should not have been called, but focusedID = %d", runner.focusedID)
		}
	})

	t.Run("runner focus error", func(t *testing.T) {
		runner := &mockRunner{focusErr: errors.New("focus failed")}
		err := FocusWindow(runner, testWindows[0], testWindows)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestFindFocusedWindow(t *testing.T) {
	t.Run("returns index of focused window", func(t *testing.T) {
		windows := []Window{
			{ID: 1, Focused: false},
			{ID: 2, Focused: true},
			{ID: 3, Focused: false},
		}
		idx := FindFocusedWindow(windows)
		if idx != 1 {
			t.Errorf("got %d, want 1", idx)
		}
	})

	t.Run("first window focused", func(t *testing.T) {
		windows := []Window{
			{ID: 1, Focused: true},
			{ID: 2, Focused: false},
		}
		idx := FindFocusedWindow(windows)
		if idx != 0 {
			t.Errorf("got %d, want 0", idx)
		}
	})

	t.Run("last window focused", func(t *testing.T) {
		windows := []Window{
			{ID: 1, Focused: false},
			{ID: 2, Focused: true},
		}
		idx := FindFocusedWindow(windows)
		if idx != 1 {
			t.Errorf("got %d, want 1", idx)
		}
	})

	t.Run("none focused", func(t *testing.T) {
		windows := []Window{
			{ID: 1, Focused: false},
			{ID: 2, Focused: false},
		}
		idx := FindFocusedWindow(windows)
		if idx != -1 {
			t.Errorf("got %d, want -1", idx)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		idx := FindFocusedWindow([]Window{})
		if idx != -1 {
			t.Errorf("got %d, want -1", idx)
		}
	})

	t.Run("single focused window", func(t *testing.T) {
		windows := []Window{{ID: 1, Focused: true}}
		idx := FindFocusedWindow(windows)
		if idx != 0 {
			t.Errorf("got %d, want 0", idx)
		}
	})

	t.Run("single unfocused window", func(t *testing.T) {
		windows := []Window{{ID: 1, Focused: false}}
		idx := FindFocusedWindow(windows)
		if idx != -1 {
			t.Errorf("got %d, want -1", idx)
		}
	})
}

func TestGetNextWindow(t *testing.T) {
	windows := []Window{
		{ID: 1, AppID: "a"},
		{ID: 2, AppID: "a"},
		{ID: 3, AppID: "a"},
	}

	t.Run("returns next window", func(t *testing.T) {
		w := GetNextWindow(0, windows)
		if w.ID != 2 {
			t.Errorf("got ID %d, want 2", w.ID)
		}
	})

	t.Run("returns next from middle", func(t *testing.T) {
		w := GetNextWindow(1, windows)
		if w.ID != 3 {
			t.Errorf("got ID %d, want 3", w.ID)
		}
	})

	t.Run("wraps around from last index", func(t *testing.T) {
		w := GetNextWindow(2, windows)
		if w.ID != 1 {
			t.Errorf("got ID %d, want 1", w.ID)
		}
	})

	t.Run("wraps around for negative index", func(t *testing.T) {
		w := GetNextWindow(-1, windows)
		if w.ID != 1 {
			t.Errorf("got ID %d, want 1", w.ID)
		}
	})

	t.Run("wraps around for out of bounds index", func(t *testing.T) {
		w := GetNextWindow(10, windows)
		if w.ID != 1 {
			t.Errorf("got ID %d, want 1", w.ID)
		}
	})

	t.Run("single window wraps to itself", func(t *testing.T) {
		single := []Window{{ID: 5, AppID: "b"}}
		w := GetNextWindow(0, single)
		if w.ID != 5 {
			t.Errorf("got ID %d, want 5", w.ID)
		}
	})

	t.Run("empty list returns zero value", func(t *testing.T) {
		w := GetNextWindow(0, []Window{})
		if w.ID != 0 {
			t.Errorf("got ID %d, want 0 (zero value)", w.ID)
		}
	})
}

func TestLaunchOrFocus(t *testing.T) {
	t.Run("focuses when single app match is focused", func(t *testing.T) {
		runner := &mockRunner{jsonData: testWindowsJSON}
		err := LaunchOrFocus(runner, "firefox", "firefox")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// firefox is the only match and is focused (last in appIDWindows) → allWindows[0]
		if runner.focusedID != 1 {
			t.Errorf("focused ID = %d, want 1", runner.focusedID)
		}
		if runner.spawnedCmd != "" {
			t.Error("Spawn should not have been called")
		}
	})

	t.Run("focuses first appID window when none is focused", func(t *testing.T) {
		runner := &mockRunner{jsonData: testWindowsJSON}
		err := LaunchOrFocus(runner, "Alacritty", "alacritty")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// No Alacritty window is focused → GetNextWindow(-1, appIDWindows) → appIDWindows[0] = id:2
		if runner.focusedID != 2 {
			t.Errorf("focused ID = %d, want 2", runner.focusedID)
		}
	})

	t.Run("cycles to next appID window when non-last is focused", func(t *testing.T) {
		json := []byte(`[
			{"id":1,"title":"Firefox","app_id":"firefox","pid":100,"workspace_id":1,"is_focused":false,"is_floating":false},
			{"id":2,"title":"Terminal 1","app_id":"Alacritty","pid":200,"workspace_id":2,"is_focused":false,"is_floating":false},
			{"id":3,"title":"Terminal 2","app_id":"Alacritty","pid":300,"workspace_id":3,"is_focused":true,"is_floating":false},
			{"id":4,"title":"Terminal 3","app_id":"Alacritty","pid":400,"workspace_id":4,"is_focused":false,"is_floating":false}
		]`)
		runner := &mockRunner{jsonData: json}
		err := LaunchOrFocus(runner, "Alacritty", "alacritty")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// appIDWindows: [id:2, id:3(focused), id:4] → GetNextWindow(1, ...) → appIDWindows[2] = id:4
		if runner.focusedID != 4 {
			t.Errorf("focused ID = %d, want 4", runner.focusedID)
		}
	})

	t.Run("wraps to first appID window when last is focused", func(t *testing.T) {
		json := []byte(`[
			{"id":1,"title":"Firefox","app_id":"firefox","pid":100,"workspace_id":1,"is_focused":false,"is_floating":false},
			{"id":2,"title":"Terminal 1","app_id":"Alacritty","pid":200,"workspace_id":2,"is_focused":false,"is_floating":false},
			{"id":3,"title":"Terminal 2","app_id":"Alacritty","pid":300,"workspace_id":3,"is_focused":false,"is_floating":false},
			{"id":4,"title":"Terminal 3","app_id":"Alacritty","pid":400,"workspace_id":4,"is_focused":true,"is_floating":false}
		]`)
		runner := &mockRunner{jsonData: json}
		err := LaunchOrFocus(runner, "Alacritty", "alacritty")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// appIDWindows: [id:2, id:3, id:4(focused)] → GetNextWindow(2, ...) → appIDWindows[0] = id:2
		if runner.focusedID != 2 {
			t.Errorf("focused ID = %d, want 2", runner.focusedID)
		}
	})

	t.Run("spawns when app not found", func(t *testing.T) {
		runner := &mockRunner{jsonData: testWindowsJSON}
		err := LaunchOrFocus(runner, "nonexistent", "some-app --flag")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if runner.spawnedCmd != "some-app --flag" {
			t.Errorf("spawned cmd = %q, want %q", runner.spawnedCmd, "some-app --flag")
		}
		if runner.focusedID != 0 {
			t.Error("Focus should not have been called")
		}
	})

	t.Run("GetWindows error propagates", func(t *testing.T) {
		runner := &mockRunner{jsonErr: errors.New("niri not running")}
		err := LaunchOrFocus(runner, "firefox", "firefox")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("spawn error propagates", func(t *testing.T) {
		runner := &mockRunner{jsonData: []byte(`[]`), spawnErr: errors.New("spawn failed")}
		err := LaunchOrFocus(runner, "anything", "bad-cmd")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("focus error propagates", func(t *testing.T) {
		runner := &mockRunner{jsonData: testWindowsJSON, focusErr: errors.New("focus failed")}
		err := LaunchOrFocus(runner, "firefox", "firefox")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
