# sway-mcp-go

Archived compatibility surface for the older Sway-specific desktop automation
server.

The active workstation control plane now lives in
`hairglasses-studio/dotfiles` at `mcp/dotfiles-mcp/`, where Hyprland and the
current desktop automation surface are maintained. Keep this repo limited to
redirect notes and compatibility context until callers are fully migrated.

## Current Surface

- Status: archive-ready redirect
- Active merge target: `dotfiles/mcp/dotfiles-mcp`
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
make build
make check
```

## Local MCP Launch

```bash
bash ./scripts/run-sway-mcp.sh
```

Repo-local MCP configuration lives in `.mcp.json`, with curated Codex profiles in
`.codex/mcp-profile-policy.json`.

For a local checkout, the repo ships the launcher plus `.mcp.json`, so MCP
clients can attach directly without reconstructing the command manually.

## Runtime Requirements

- Wayland session with `WAYLAND_DISPLAY` set
- Sway tooling such as `swaymsg`, `grim`, `wl-copy`, and `wl-paste`
- input/screenshot helpers such as `ydotool` and ImageMagick `magick`
- The server still boots without `WAYLAND_DISPLAY` so discovery, resources,
  prompts, and health checks remain available before a live Sway session is
  attached.
