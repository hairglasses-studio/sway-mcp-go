package sway

import (
	"testing"
)

func TestModuleTools(t *testing.T) {
	m := &Module{}
	tools := m.Tools()

	expectedCount := 12
	if len(tools) != expectedCount {
		t.Errorf("expected %d tools, got %d", expectedCount, len(tools))
	}

	expectedNames := []string{
		"sway_screenshot",
		"sway_screenshot_region",
		"sway_click",
		"sway_type_text",
		"sway_key",
		"sway_scroll",
		"sway_list_windows",
		"sway_focus_window",
		"sway_move_window",
		"sway_clipboard_read",
		"sway_clipboard_write",
		"sway_get_outputs",
	}

	for i, name := range expectedNames {
		if tools[i].Tool.Name != name {
			t.Errorf("expected tool %d name to be %s, got %s", i, name, tools[i].Tool.Name)
		}
	}
}
