# sway-mcp-go — Agent Instructions

> Canonical instructions: AGENTS.md

Use this file as the canonical multi-provider instruction surface for this repo. `CLAUDE.md`, `GEMINI.md`, and `.github/copilot-instructions.md` are compatibility mirrors.

## Purpose

- Public Sway-focused MCP server retained as its own repo even though the main
  workstation control plane has moved to Hyprland.

## Build & Test

- `go build ./cmd/sway-mcp`
- `go test ./... -count=1`
- `go vet ./...`

## Architecture

- `cmd/sway-mcp` boots the stdio MCP server and registers the Sway module.
- `internal/sway/sway.go` exposes the desktop automation surface.
- Runtime behavior depends on Sway/Wayland tools such as `swaymsg`, `grim`,
  `wl-copy`, `wl-paste`, `ydotool`, and `magick`.

## Working Rules

- Preserve the public `sway_` tool names unless downstream callers are updated
  in the same change.
- Keep runtime assumptions documented in `README.md` instead of burying them in
  code comments.
- Keep [AGENTS.md](AGENTS.md), [CLAUDE.md](CLAUDE.md), [GEMINI.md](GEMINI.md),
  and `.github/copilot-instructions.md` aligned when instructions change.
- Save reusable research findings back to `~/hairglasses-studio/docs/` instead
  of adding duplicate local research docs.
