# sway-mcp-go — Agent Instructions

> Canonical instructions: AGENTS.md

Use this file as the canonical multi-provider instruction surface for this repo. `CLAUDE.md`, `GEMINI.md`, and `.github/copilot-instructions.md` are compatibility mirrors.

## Purpose

- Archived compatibility surface for the older Sway-focused MCP server.
- Active desktop automation work now lives in `dotfiles/mcp/dotfiles-mcp`.
- Treat this repo as an archive-ready redirect, not an active feature surface.

## Build & Test

- `go build ./cmd/sway-mcp`
- `go test ./... -count=1`
- `go vet ./...`
- `bash ./scripts/run-sway-mcp.sh`

## Architecture

- `cmd/sway-mcp` boots the stdio MCP server plus the discovery/resources/prompts
  contract surface.
- `internal/sway/sway.go` exposes the desktop automation surface.
- `internal/sway/contract.go` defines the discovery tools, resources, prompts,
  and server health surface.
- `scripts/run-sway-mcp.sh` is the portable repo-local launcher referenced by
  `.mcp.json`.
- Runtime behavior depends on Sway/Wayland tools such as `swaymsg`, `grim`,
  `wl-copy`, `wl-paste`, `ydotool`, and `magick`.

## Working Rules

- Preserve the public `sway_` tool names unless downstream callers are updated
  in the same change.
- Do not add new workstation automation here; land new desktop work in
  `dotfiles/mcp/dotfiles-mcp`.
- Keep runtime assumptions documented in `README.md` instead of burying them in
  code comments.
- Keep [AGENTS.md](AGENTS.md), [CLAUDE.md](CLAUDE.md), [GEMINI.md](GEMINI.md),
  and `.github/copilot-instructions.md` aligned when instructions change.
- Save reusable research findings back to `~/hairglasses-studio/docs/` instead
  of adding duplicate local research docs.
