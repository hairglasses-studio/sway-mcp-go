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
- discovery and contract: `sway_tool_catalog`, `sway_tool_search`,
  `sway_tool_schema`, `sway_tool_stats`, `sway_server_health`
- resources: `sway://server/overview`, `sway://runtime/requirements`
- prompts: `sway_start_triage`

## Build & Test

```bash
go build ./cmd/sway-mcp
go test ./... -count=1
go vet ./...
```

## Local MCP Launch

```bash
bash ./scripts/run-sway-mcp.sh
```

Repo-local MCP configuration lives in `.mcp.json`, with curated Codex profiles in
`.codex/mcp-profile-policy.json`.

## Runtime Requirements

- Wayland session with `WAYLAND_DISPLAY` set
- Sway tooling such as `swaymsg`, `grim`, `wl-copy`, and `wl-paste`
- input/screenshot helpers such as `ydotool` and ImageMagick `magick`
- The server still boots without `WAYLAND_DISPLAY` so discovery, resources,
  prompts, and health checks remain available before a live Sway session is
  attached.
