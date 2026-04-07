package sway

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/registry"
)

const (
	MaxDim = 1568
)

var buttonMap = map[string]string{
	"left":   "0xC0",
	"right":  "0xC1",
	"middle": "0xC2",
	"up":     "0xC3",
	"down":   "0xC4",
}

// --- Types ---

type ScreenshotInput struct {
	Output string `json:"output,omitempty" jsonschema:"description=Output name (e.g. DP-1, DP-2). Omit for all outputs."`
}

type ScreenshotRegionInput struct {
	X      int `json:"x" jsonschema:"required,description=X coordinate"`
	Y      int `json:"y" jsonschema:"required,description=Y coordinate"`
	Width  int `json:"width" jsonschema:"required,description=Width in pixels"`
	Height int `json:"height" jsonschema:"required,description=Height in pixels"`
}

type ClickInput struct {
	X      int    `json:"x" jsonschema:"required,description=X pixel coordinate"`
	Y      int    `json:"y" jsonschema:"required,description=Y pixel coordinate"`
	Button string `json:"button,omitempty" jsonschema:"enum=left,enum=right,enum=middle,description=Mouse button (default: left)"`
	Clicks int    `json:"clicks,omitempty" jsonschema:"description=Number of clicks (default: 1, use 2 for double-click)"`
}

type TypeTextInput struct {
	Text string `json:"text" jsonschema:"required,description=Text to type"`
}

type KeyInput struct {
	Combo string `json:"combo" jsonschema:"required,description=Key combo, e.g. 'ctrl+shift+t', 'super+d', 'Return'"`
}

type ScrollInput struct {
	X         int    `json:"x" jsonschema:"required,description=X coordinate"`
	Y         int    `json:"y" jsonschema:"required,description=Y coordinate"`
	Direction string `json:"direction" jsonschema:"enum=up,enum=down,required,description=Scroll direction"`
	Amount    int    `json:"amount,omitempty" jsonschema:"description=Number of scroll steps (default: 3)"`
}

type FocusWindowInput struct {
	ConID int    `json:"con_id,omitempty" jsonschema:"description=Container ID from sway_list_windows"`
	AppID string `json:"app_id,omitempty" jsonschema:"description=Application ID (e.g. com.mitchellh.ghostty)"`
}

type MoveWindowInput struct {
	ConID  int    `json:"con_id,omitempty"`
	AppID  string `json:"app_id,omitempty"`
	X      *int   `json:"x,omitempty" jsonschema:"description=New X position"`
	Y      *int   `json:"y,omitempty" jsonschema:"description=New Y position"`
	Width  *int   `json:"width,omitempty" jsonschema:"description=New width"`
	Height *int   `json:"height,omitempty" jsonschema:"description=New height"`
}

type ClipboardWriteInput struct {
	Text string `json:"text" jsonschema:"required,description=Text to copy to clipboard"`
}

type Window struct {
	ConID     int    `json:"con_id"`
	AppID     string `json:"app_id"`
	Title     string `json:"title"`
	Rect      Rect   `json:"rect"`
	Focused   bool   `json:"focused"`
	Workspace string `json:"workspace"`
	Floating  bool   `json:"floating"`
	Urgent    bool   `json:"urgent"`
}

type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// --- Module ---

type Module struct{}

func (m *Module) Name() string { return "sway" }
func (m *Module) Description() string { return "Sway/Wayland desktop management" }

func (m *Module) Tools() []registry.ToolDefinition {
	return []registry.ToolDefinition{
		handler.TypedHandler[ScreenshotInput, string](
			"sway_screenshot",
			"Take a screenshot of the full desktop or a specific output. Returns a base64-encoded PNG.",
			m.Screenshot,
		),
		handler.TypedHandler[ScreenshotRegionInput, string](
			"sway_screenshot_region",
			"Take a screenshot of a specific rectangular region. Returns a base64-encoded PNG.",
			m.ScreenshotRegion,
		),
		handler.TypedHandler[ClickInput, string](
			"sway_click",
			"Click at absolute pixel coordinates.",
			m.Click,
		),
		handler.TypedHandler[TypeTextInput, string](
			"sway_type_text",
			"Type text into the currently focused window.",
			m.TypeText,
		),
		handler.TypedHandler[KeyInput, string](
			"sway_key",
			"Send a keyboard shortcut.",
			m.Key,
		),
		handler.TypedHandler[ScrollInput, string](
			"sway_scroll",
			"Scroll at a position.",
			m.Scroll,
		),
		handler.TypedHandler[struct{}, []Window](
			"sway_list_windows",
			"List all windows.",
			m.ListWindows,
		),
		handler.TypedHandler[FocusWindowInput, string](
			"sway_focus_window",
			"Focus a window by con_id or app_id.",
			m.FocusWindow,
		),
		handler.TypedHandler[MoveWindowInput, string](
			"sway_move_window",
			"Move and/or resize a window.",
			m.MoveWindow,
		),
		handler.TypedHandler[struct{}, string](
			"sway_clipboard_read",
			"Read the current Wayland clipboard contents.",
			m.ClipboardRead,
		),
		handler.TypedHandler[ClipboardWriteInput, string](
			"sway_clipboard_write",
			"Write text to the Wayland clipboard.",
			m.ClipboardWrite,
		),
		handler.TypedHandler[struct{}, interface{}](
			"sway_get_outputs",
			"List all monitors/outputs.",
			m.GetOutputs,
		),
	}
}

// --- Implementation ---

func (m *Module) Screenshot(ctx context.Context, in ScreenshotInput) (string, error) {
	args := []string{}
	if in.Output != "" {
		args = append(args, "-o", in.Output)
	}
	return m.captureAndScale(ctx, args)
}

func (m *Module) ScreenshotRegion(ctx context.Context, in ScreenshotRegionInput) (string, error) {
	args := []string{"-g", fmt.Sprintf("%d,%d %dx%d", in.X, in.Y, in.Width, in.Height)}
	return m.captureAndScale(ctx, args)
}

func (m *Module) captureAndScale(ctx context.Context, grimArgs []string) (string, error) {
	raw := filepath.Join(os.TempDir(), "sway-mcp-raw.png")
	scaled := filepath.Join(os.TempDir(), "sway-mcp-scaled.png")
	defer os.Remove(raw)
	defer os.Remove(scaled)

	args := append(grimArgs, raw)
	if err := exec.CommandContext(ctx, "grim", args...).Run(); err != nil {
		return "", fmt.Errorf("grim failed: %w", err)
	}

	magickArgs := []string{"convert", raw, "-resize", fmt.Sprintf("%dx%d\\>", MaxDim, MaxDim), scaled}
	if err := exec.CommandContext(ctx, "magick", magickArgs...).Run(); err != nil {
		return "", fmt.Errorf("magick failed: %w", err)
	}

	data, err := os.ReadFile(scaled)
	if err != nil {
		return "", fmt.Errorf("read scaled file: %w", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func (m *Module) Click(ctx context.Context, in ClickInput) (string, error) {
	if err := exec.CommandContext(ctx, "ydotool", "mousemove", "--absolute", "-x", strconv.Itoa(in.X), "-y", strconv.Itoa(in.Y)).Run(); err != nil {
		return "", fmt.Errorf("ydotool mousemove failed: %w", err)
	}

	button := in.Button
	if button == "" {
		button = "left"
	}
	btn, ok := buttonMap[button]
	if !ok {
		btn = buttonMap["left"]
	}

	clicks := in.Clicks
	if clicks == 0 {
		clicks = 1
	}

	for i := 0; i < clicks; i++ {
		if err := exec.CommandContext(ctx, "ydotool", "click", btn).Run(); err != nil {
			return "", fmt.Errorf("ydotool click failed: %w", err)
		}
	}
	return fmt.Sprintf("Clicked %s at %d,%d", button, in.X, in.Y), nil
}

func (m *Module) TypeText(ctx context.Context, in TypeTextInput) (string, error) {
	if err := exec.CommandContext(ctx, "wtype", in.Text).Run(); err != nil {
		return "", fmt.Errorf("wtype failed: %w", err)
	}
	return fmt.Sprintf("Typed %d characters", len(in.Text)), nil
}

func (m *Module) Key(ctx context.Context, in KeyInput) (string, error) {
	parts := strings.Split(strings.ToLower(in.Combo), "+")
	keyName := parts[len(parts)-1]
	mods := parts[:len(parts)-1]

	args := []string{}
	for _, mod := range mods {
		args = append(args, "-M", strings.TrimSpace(mod))
	}
	args = append(args, "-k", strings.TrimSpace(keyName))
	for i := len(mods) - 1; i >= 0; i-- {
		args = append(args, "-m", strings.TrimSpace(mods[i]))
	}

	if err := exec.CommandContext(ctx, "wtype", args...).Run(); err != nil {
		return "", fmt.Errorf("wtype failed: %w", err)
	}
	return fmt.Sprintf("Sent %s", in.Combo), nil
}

func (m *Module) Scroll(ctx context.Context, in ScrollInput) (string, error) {
	if err := exec.CommandContext(ctx, "ydotool", "mousemove", "--absolute", "-x", strconv.Itoa(in.X), "-y", strconv.Itoa(in.Y)).Run(); err != nil {
		return "", fmt.Errorf("ydotool mousemove failed: %w", err)
	}

	btn, ok := buttonMap[in.Direction]
	if !ok {
		return "", fmt.Errorf("invalid scroll direction: %s", in.Direction)
	}

	amount := in.Amount
	if amount == 0 {
		amount = 3
	}

	for i := 0; i < amount; i++ {
		if err := exec.CommandContext(ctx, "ydotool", "click", btn).Run(); err != nil {
			return "", fmt.Errorf("ydotool click failed: %w", err)
		}
	}
	return fmt.Sprintf("Scrolled %s %d at %d,%d", in.Direction, amount, in.X, in.Y), nil
}

func (m *Module) ListWindows(ctx context.Context, _ struct{}) ([]Window, error) {
	out, err := exec.CommandContext(ctx, "swaymsg", "-t", "get_tree").Output()
	if err != nil {
		return nil, fmt.Errorf("swaymsg get_tree failed: %w", err)
	}

	var tree map[string]interface{}
	if err := json.Unmarshal(out, &tree); err != nil {
		return nil, fmt.Errorf("unmarshal sway tree: %w", err)
	}

	return flattenTree(tree, ""), nil
}

func flattenTree(node map[string]interface{}, workspace string) []Window {
	results := []Window{}
	ws := workspace
	if node["type"] == "workspace" {
		if name, ok := node["name"].(string); ok {
			ws = name
		}
	}

	appID, _ := node["app_id"].(string)
	if appID == "" {
		if props, ok := node["window_properties"].(map[string]interface{}); ok {
			appID, _ = props["class"].(string)
		}
	}

	pid, _ := node["pid"].(float64)

	if appID != "" && pid != 0 {
		rectMap, _ := node["rect"].(map[string]interface{})
		results = append(results, Window{
			ConID: int(node["id"].(float64)),
			AppID: appID,
			Title: node["name"].(string),
			Rect: Rect{
				X:      int(rectMap["x"].(float64)),
				Y:      int(rectMap["y"].(float64)),
				Width:  int(rectMap["width"].(float64)),
				Height: int(rectMap["height"].(float64)),
			},
			Focused:   node["focused"].(bool),
			Workspace: ws,
			Floating:  node["type"] == "floating_con",
			Urgent:    node["urgent"].(bool),
		})
	}

	if nodes, ok := node["nodes"].([]interface{}); ok {
		for _, n := range nodes {
			results = append(results, flattenTree(n.(map[string]interface{}), ws)...)
		}
	}
	if fnodes, ok := node["floating_nodes"].([]interface{}); ok {
		for _, n := range fnodes {
			results = append(results, flattenTree(n.(map[string]interface{}), ws)...)
		}
	}

	return results
}

func (m *Module) FocusWindow(ctx context.Context, in FocusWindowInput) (string, error) {
	var selector string
	if in.ConID != 0 {
		selector = fmt.Sprintf("[con_id=%d]", in.ConID)
	} else if in.AppID != "" {
		selector = fmt.Sprintf("[app_id=\"%s\"]", in.AppID)
	} else {
		return "", fmt.Errorf("provide con_id or app_id")
	}

	if err := exec.CommandContext(ctx, "swaymsg", selector, "focus").Run(); err != nil {
		return "", fmt.Errorf("swaymsg focus failed: %w", err)
	}
	return "Focused", nil
}

func (m *Module) MoveWindow(ctx context.Context, in MoveWindowInput) (string, error) {
	var selector string
	if in.ConID != 0 {
		selector = fmt.Sprintf("[con_id=%d]", in.ConID)
	} else if in.AppID != "" {
		selector = fmt.Sprintf("[app_id=\"%s\"]", in.AppID)
	} else {
		return "", fmt.Errorf("provide con_id or app_id")
	}

	if in.X != nil && in.Y != nil {
		if err := exec.CommandContext(ctx, "swaymsg", selector, "move", "position", strconv.Itoa(*in.X), strconv.Itoa(*in.Y)).Run(); err != nil {
			return "", fmt.Errorf("swaymsg move position failed: %w", err)
		}
	}
	if in.Width != nil && in.Height != nil {
		if err := exec.CommandContext(ctx, "swaymsg", selector, "resize", "set", strconv.Itoa(*in.Width), strconv.Itoa(*in.Height)).Run(); err != nil {
			return "", fmt.Errorf("swaymsg resize set failed: %w", err)
		}
	}
	return "Done", nil
}

func (m *Module) ClipboardRead(ctx context.Context, _ struct{}) (string, error) {
	out, err := exec.CommandContext(ctx, "wl-paste", "-n").Output()
	if err != nil {
		return "", fmt.Errorf("wl-paste failed: %w", err)
	}
	return string(out), nil
}

func (m *Module) ClipboardWrite(ctx context.Context, in ClipboardWriteInput) (string, error) {
	cmd := exec.CommandContext(ctx, "wl-copy")
	cmd.Stdin = strings.NewReader(in.Text)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("wl-copy failed: %w", err)
	}
	return "Copied to clipboard", nil
}

func (m *Module) GetOutputs(ctx context.Context, _ struct{}) (interface{}, error) {
	out, err := exec.CommandContext(ctx, "swaymsg", "-t", "get_outputs").Output()
	if err != nil {
		return nil, fmt.Errorf("swaymsg get_outputs failed: %w", err)
	}

	var results interface{}
	if err := json.Unmarshal(out, &results); err != nil {
		return nil, fmt.Errorf("unmarshal sway outputs: %w", err)
	}
	return results, nil
}
