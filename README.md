# sway-mcp-go

`sway-mcp-go` is a standalone Go MCP server for Sway and Wayland desktop
automation.

The main `hairglasses-studio` desktop stack now targets Hyprland, but this repo
is retained as a separate Sway-focused surface.

## Current Surface

- screenshots: `sway_screenshot`, `sway_screenshot_region`
- input: `sway_click`, `sway_type_text`, `sway_key`, `sway_scroll`
- window management: `sway_list_windows`, `sway_focus_window`,
  `sway_move_window`
- clipboard and display inspection: `sway_clipboard_read`,
  `sway_clipboard_write`, `sway_get_outputs`

## Build & Test

```bash
go build ./cmd/sway-mcp
go test ./... -count=1
go vet ./...
```

## Runtime Requirements

- Wayland session with `WAYLAND_DISPLAY` set
- Sway tooling such as `swaymsg`, `grim`, `wl-copy`, and `wl-paste`
- input/screenshot helpers such as `ydotool` and ImageMagick `magick`
