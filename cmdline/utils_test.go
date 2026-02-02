package cmdline

import (
	"strings"
	"testing"

	nirilof "github.com/azr4e1/niri-lof"
)

func TestValidateFormat(t *testing.T) {
	t.Run("empty string returns empty slice", func(t *testing.T) {
		opts, err := ValidateFormat("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(opts) != 0 {
			t.Errorf("expected empty slice, got %v", opts)
		}
	})

	t.Run("whitespace only returns empty slice", func(t *testing.T) {
		opts, err := ValidateFormat("   ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(opts) != 0 {
			t.Errorf("expected empty slice, got %v", opts)
		}
	})

	t.Run("single valid option", func(t *testing.T) {
		opts, err := ValidateFormat("ID")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(opts) != 1 || opts[0] != "ID" {
			t.Errorf("expected [ID], got %v", opts)
		}
	})

	t.Run("multiple valid options", func(t *testing.T) {
		opts, err := ValidateFormat("ID,APPID,TITLE")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := []string{"ID", "APPID", "TITLE"}
		if len(opts) != len(expected) {
			t.Fatalf("expected %d options, got %d", len(expected), len(opts))
		}
		for i, e := range expected {
			if opts[i] != e {
				t.Errorf("opts[%d] = %q, want %q", i, opts[i], e)
			}
		}
	})

	t.Run("all valid options", func(t *testing.T) {
		opts, err := ValidateFormat("ID,APPID,TITLE,FOCUSED,ISFLOATING,PID,WORKSPACEID")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(opts) != 7 {
			t.Errorf("expected 7 options, got %d", len(opts))
		}
	})

	t.Run("invalid option returns error", func(t *testing.T) {
		_, err := ValidateFormat("INVALID")
		if err == nil {
			t.Fatal("expected error for invalid option, got nil")
		}
		if !strings.Contains(err.Error(), "INVALID") {
			t.Errorf("error should mention the invalid option, got: %v", err)
		}
	})

	t.Run("mixed valid and invalid returns error", func(t *testing.T) {
		_, err := ValidateFormat("ID,BOGUS,TITLE")
		if err == nil {
			t.Fatal("expected error for invalid option, got nil")
		}
		if !strings.Contains(err.Error(), "BOGUS") {
			t.Errorf("error should mention the invalid option, got: %v", err)
		}
	})

	t.Run("duplicate valid options are accepted", func(t *testing.T) {
		opts, err := ValidateFormat("ID,ID")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(opts) != 2 || opts[0] != "ID" || opts[1] != "ID" {
			t.Errorf("expected [ID ID], got %v", opts)
		}
	})

	t.Run("whitespace-padded input is trimmed", func(t *testing.T) {
		opts, err := ValidateFormat("  ID,APPID  ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(opts) != 2 || opts[0] != "ID" || opts[1] != "APPID" {
			t.Errorf("expected [ID APPID], got %v", opts)
		}
	})

	t.Run("lowercase option is invalid", func(t *testing.T) {
		_, err := ValidateFormat("id")
		if err == nil {
			t.Fatal("expected error for lowercase option, got nil")
		}
	})
}

var testWindows = []nirilof.Window{
	{ID: 1, Title: "Firefox", AppID: "firefox", PID: 100, WorkspaceID: 1, Focused: true, IsFloating: false},
	{ID: 2, Title: "Terminal", AppID: "Alacritty", PID: 200, WorkspaceID: 2, Focused: false, IsFloating: false},
	{ID: 3, Title: "Floating Term", AppID: "Alacritty", PID: 300, WorkspaceID: 3, Focused: false, IsFloating: true},
}

func TestFormatWindows(t *testing.T) {
	t.Run("single column", func(t *testing.T) {
		result := FormatWindows(testWindows, []string{"APPID"})
		lines := strings.Split(result, "\n")
		if len(lines) != 4 { // header + 3 windows
			t.Fatalf("expected 4 lines, got %d", len(lines))
		}
		if strings.TrimSpace(lines[0]) != "APPID" {
			t.Errorf("header = %q, want %q", strings.TrimSpace(lines[0]), "APPID")
		}
		if strings.TrimSpace(lines[1]) != "firefox" {
			t.Errorf("row 1 = %q, want %q", strings.TrimSpace(lines[1]), "firefox")
		}
		if strings.TrimSpace(lines[2]) != "Alacritty" {
			t.Errorf("row 2 = %q, want %q", strings.TrimSpace(lines[2]), "Alacritty")
		}
	})

	t.Run("multiple columns", func(t *testing.T) {
		result := FormatWindows(testWindows, []string{"ID", "APPID", "FOCUSED"})
		lines := strings.Split(result, "\n")
		if len(lines) != 4 {
			t.Fatalf("expected 4 lines, got %d", len(lines))
		}
		// Header row should contain all three column names
		header := lines[0]
		for _, col := range []string{"ID", "APPID", "FOCUSED"} {
			if !strings.Contains(header, col) {
				t.Errorf("header missing column %q: %q", col, header)
			}
		}
		// First data row should have firefox's data
		if !strings.Contains(lines[1], "1") || !strings.Contains(lines[1], "firefox") || !strings.Contains(lines[1], "true") {
			t.Errorf("row 1 unexpected: %q", lines[1])
		}
	})

	t.Run("empty window list has header only", func(t *testing.T) {
		result := FormatWindows([]nirilof.Window{}, []string{"ID", "APPID"})
		lines := strings.Split(result, "\n")
		if len(lines) != 1 {
			t.Fatalf("expected 1 line (header only), got %d", len(lines))
		}
		if !strings.Contains(lines[0], "ID") || !strings.Contains(lines[0], "APPID") {
			t.Errorf("header unexpected: %q", lines[0])
		}
	})

	t.Run("default columns when options is empty", func(t *testing.T) {
		result := FormatWindows(testWindows, []string{})
		lines := strings.Split(result, "\n")
		if len(lines) != 4 {
			t.Fatalf("expected 4 lines, got %d", len(lines))
		}
		header := lines[0]
		for _, col := range []string{"ID", "APPID", "TITLE", "FOCUSED", "ISFLOATING", "PID", "WORKSPACEID"} {
			if !strings.Contains(header, col) {
				t.Errorf("default header missing column %q: %q", col, header)
			}
		}
	})

	t.Run("columns are padded to align", func(t *testing.T) {
		result := FormatWindows(testWindows, []string{"APPID"})
		lines := strings.Split(result, "\n")
		// "Alacritty" (9 chars) is the longest value, so all entries
		// should be padded to at least 9 chars
		for _, line := range lines {
			// Each line is a single column, so just check it's padded
			if len(line) < len("Alacritty") {
				t.Errorf("line %q shorter than longest value", line)
			}
		}
	})

	t.Run("columns separated by double space", func(t *testing.T) {
		result := FormatWindows(testWindows, []string{"ID", "APPID"})
		lines := strings.Split(result, "\n")
		for i, line := range lines {
			if !strings.Contains(line, "  ") {
				t.Errorf("line %d missing double space separator: %q", i, line)
			}
		}
	})

	t.Run("boolean fields render correctly", func(t *testing.T) {
		result := FormatWindows(testWindows, []string{"FOCUSED", "ISFLOATING"})
		lines := strings.Split(result, "\n")
		// line[1] is firefox: focused=true, floating=false
		if !strings.Contains(lines[1], "true") || !strings.Contains(lines[1], "false") {
			t.Errorf("row 1 unexpected: %q", lines[1])
		}
		// line[3] is floating term: focused=false, floating=true
		if !strings.Contains(lines[3], "false") || !strings.Contains(lines[3], "true") {
			t.Errorf("row 3 unexpected: %q", lines[3])
		}
	})

	t.Run("single window", func(t *testing.T) {
		single := []nirilof.Window{
			{ID: 42, Title: "Solo", AppID: "solo", PID: 999, WorkspaceID: 1},
		}
		result := FormatWindows(single, []string{"ID", "TITLE"})
		lines := strings.Split(result, "\n")
		if len(lines) != 2 {
			t.Fatalf("expected 2 lines, got %d", len(lines))
		}
		if !strings.Contains(lines[1], "42") || !strings.Contains(lines[1], "Solo") {
			t.Errorf("row 1 unexpected: %q", lines[1])
		}
	})

	t.Run("all column types render", func(t *testing.T) {
		w := []nirilof.Window{
			{ID: 5, Title: "Test", AppID: "test", PID: 555, WorkspaceID: 2, Focused: true, IsFloating: true},
		}
		result := FormatWindows(w, []string{"ID", "APPID", "TITLE", "FOCUSED", "ISFLOATING", "PID", "WORKSPACEID"})
		lines := strings.Split(result, "\n")
		if len(lines) != 2 {
			t.Fatalf("expected 2 lines, got %d", len(lines))
		}
		row := lines[1]
		for _, expected := range []string{"5", "test", "Test", "true", "true", "555", "2"} {
			if !strings.Contains(row, expected) {
				t.Errorf("row missing %q: %q", expected, row)
			}
		}
	})
}
